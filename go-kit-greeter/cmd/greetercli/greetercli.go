package main

import (
	"io"
	"time"

	"github.com/antklim/go-microservices/go-kit-greeter/pkg/greeterendpoint"
	"github.com/antklim/go-microservices/go-kit-greeter/pkg/greeterservice"
	"github.com/antklim/go-microservices/go-kit-greeter/pkg/greetertransport"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	consulapi "github.com/hashicorp/consul/api"
)

// New returns a service that's load-balanced over instances of greeterservice
// found in the provided Consul server. The mechanism of looking up greeterservice
// instances in Consul is hard-coded into the client.
func New(consulAddr string, logger log.Logger) (greeterservice.Service, error) {
	apiclient, err := consulapi.NewClient(&consulapi.Config{
		Address: consulAddr,
	})
	if err != nil {
		return nil, err
	}

	var (
		consulService = "go.kit.srv.greeter"
		consulTags    = []string{}
		passingOnly   = true
		retryMax      = 3
		retryTimeout  = 500 * time.Millisecond
	)

	var (
		sdclient  = consul.NewClient(apiclient)
		instancer = consul.NewInstancer(sdclient, logger, consulService, consulTags, passingOnly)
		endpoints greeterendpoint.Endpoints
	)
	{
		factory := factoryFor(greeterendpoint.MakeHealthEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.HealthEndpoint = retry
	}
	{
		factory := factoryFor(greeterendpoint.MakeGreetingEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.GreetingEndpoint = retry
	}

	return endpoints, nil
}

func factoryFor(makeEndpoint func(greeterservice.Service) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		service, err := greetertransport.MakeHTTPClientEndpoints(instance)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}

func main() {

}
