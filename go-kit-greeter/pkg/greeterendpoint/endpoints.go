package greeterendpoint

import (
	"context"

	"github.com/go-kit/kit/log"

	"github.com/antklim/go-microservices/go-kit-greeter/pkg/greeterservice"
	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that compose a greeter service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	HealthEndpoint   endpoint.Endpoint // used by Consul for the healthcheck
	GreetingEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns service Endoints, and wires in all the provided
// middlewares.
func MakeServerEndpoints(s greeterservice.Service, logger log.Logger) Endpoints {
	var healthEndpoint endpoint.Endpoint
	{
		healthEndpoint = MakeHealthEndpoint(s)
		healthEndpoint = LoggingMiddleware(log.With(logger, "method", "Health"))(healthEndpoint)
	}

	var greetingEndpoint endpoint.Endpoint
	{
		greetingEndpoint = MakeGreetingEndpoint(s)
		greetingEndpoint = LoggingMiddleware(log.With(logger, "method", "Greeting"))(greetingEndpoint)
	}

	return Endpoints{
		HealthEndpoint:   healthEndpoint,
		GreetingEndpoint: greetingEndpoint,
	}
}

// MakeHealthEndpoint constructs a Health endpoint wrapping the service.
func MakeHealthEndpoint(s greeterservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		healthy := s.Health()
		return HealthResponse{Healthy: healthy}, nil
	}
}

// MakeGreetingEndpoint constructs a Greeter endpoint wrapping the service.
func MakeGreetingEndpoint(s greeterservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GreetingRequest)
		greeting := s.Greeting(req.Name)
		return GreetingResponse{Greeting: greeting}, nil
	}
}

// Health implements greeterservice.Service.
// func (e Endpoints) Health(ctx context.Context, request interface{}) (bool, error) {
// 	response, err := e.HealthEndpoint(ctx, request)
// 	if err != nil {
// 		return false, err
// 	}
// 	resp := response.(HealthResponse)
// 	return resp.Health, err
// }

// // Greeting implements greeterservice.Service.
// func (e Endpoints) Greeting(ctx context.Context, name string) (string, error) {
// 	response, err := e.GreetingEndpoint(ctx, name)
// 	if err != nil {
// 		return "", err
// 	}
// 	resp := response.(GreetingResponse)
// 	return resp.Greeting, err
// }

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type Failer interface {
	Failed() error
}

// HealthRequest collects the request parameters for the Health method.
type HealthRequest struct{}

// HealthResponse collects the response values for the Health method.
type HealthResponse struct {
	Healthy bool  `json:"healthy,omitempty"`
	Err     error `json:"err,omitempty"`
}

// Failed implements Failer.
func (r HealthResponse) Failed() error { return r.Err }

// GreetingRequest collects the request parameters for the Greeting method.
type GreetingRequest struct {
	Name string `json:"name,omitempty"`
}

// GreetingResponse collects the response values for the Greeting method.
type GreetingResponse struct {
	Greeting string `json:"greeting,omitempty"`
	Err      error  `json:"err,omitempty"`
}

// Failed implements Failer.
func (r GreetingResponse) Failed() error { return r.Err }
