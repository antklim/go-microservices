package gokitgreeter

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that compose a greeter service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	GetHealthEndpoint   endpoint.Endpoint // used by Consul for the healthcheck
	GetGreetingEndpoint endpoint.Endpoint
}

func MakeServiceEndpoints(s Service) Endpoints {
	return Endpoints{
		GetHealthEndpoint:   MakeGetHealthEndpoint(s),
		GetGreetingEndpoint: MakeGetGreetingEndpoint(s),
	}
}

func MakeGetHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		health, err := s.GetHealth(ctx)
		return getHealthResponse{Health: health, Err: err}, nil
	}
}

func MakeGetGreetingEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getGreetingRequest)
		greeting, err := s.GetGreeting(ctx, req.Name)
		return getGreetingResponse{Greeting: greeting, Err: err}, nil
	}
}

type getHealthRequest struct{}

type getHealthResponse struct {
	Health bool  `json:"health,omitempty"`
	Err    error `json:"err,omitempty"`
}

func (r getHealthResponse) error() error { return r.Err }

type getGreetingRequest struct {
	Name string
}

type getGreetingResponse struct {
	Greeting string `json:"greeting,omitempty"`
	Err      error  `json:"err,omitempty"`
}

func (r getGreetingResponse) error() error { return r.Err }
