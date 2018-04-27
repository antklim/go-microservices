package greeterendpoint

import (
	"net/http"

	ocontext "golang.org/x/net/context"

	"github.com/NYTimes/gizmo/server"
	"github.com/antklim/go-microservices/gizmo-greeter/pkg/greeterservice"
)

// Endpoints collects all of the endpoints that compose a greeter service.
type Endpoints struct {
	HealthEndpoint   server.JSONContextEndpoint
	GreetingEndpoint server.JSONContextEndpoint
}

// MakeServerEndpoints returns service Endoints
func MakeServerEndpoints(s greeterservice.Service) Endpoints {
	healthEndpoint := MakeHealthEndpoint(s)
	greetingEndpoint := MakeGreetingEndpoint(s)

	return Endpoints{
		HealthEndpoint:   healthEndpoint,
		GreetingEndpoint: greetingEndpoint,
	}
}

// MakeHealthEndpoint constructs a Health endpoint.
func MakeHealthEndpoint(s greeterservice.Service) server.JSONContextEndpoint {
	return func(ctx ocontext.Context, r *http.Request) (int, interface{}, error) {
		healthy := s.Health()
		return http.StatusOK, HealthResponse{Healthy: healthy}, nil
	}
}

// MakeGreetingEndpoint constructs a Greeting endpoint.
func MakeGreetingEndpoint(s greeterservice.Service) server.JSONContextEndpoint {
	return func(ctx ocontext.Context, r *http.Request) (int, interface{}, error) {
		// TODO - get name parameter from the request query
		greeting := s.Greeting("BOB")
		return http.StatusOK, GreetingResponse{Greeting: greeting}, nil
	}
}

// HealthRequest collects the request parameters for the Health method.
type HealthRequest struct{}

// HealthResponse collects the response values for the Health method.
type HealthResponse struct {
	Healthy bool `json:"healthy,omitempty"`
}

// GreetingRequest collects the request parameters for the Greeting method.
type GreetingRequest struct {
	Name string `json:"name,omitempty"`
}

// GreetingResponse collects the response values for the Greeting method.
type GreetingResponse struct {
	Greeting string `json:"greeting,omitempty"`
}
