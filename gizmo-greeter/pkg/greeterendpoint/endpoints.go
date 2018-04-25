package greeterendpoint

import (
	"net/http"

	"github.com/NYTimes/gizmo/server"
)

// Endpoints collects all of the endpoints that compose a greeter service.
type Endpoints struct {
	HealthEndpoint   server.JSONEndpoint
	GreetingEndpoint server.JSONEndpoint
}

// MakeServerEndpoints returns service Endoints
func MakeServerEndpoints() Endpoints {
	healthEndpoint := MakeHealthEndpoint()
	greetingEndpoint := MakeGreetingEndpoint()

	return Endpoints{
		HealthEndpoint:   healthEndpoint,
		GreetingEndpoint: greetingEndpoint,
	}
}

// MakeHealthEndpoint constructs a Health endpoint.
func MakeHealthEndpoint() server.JSONEndpoint {
	return func(r *http.Request) (int, interface{}, error) {
		// TODO - add real service call to get health state
		healthy := true
		return http.StatusOK, HealthResponse{Healthy: healthy}, nil
	}
}

// MakeGreetingEndpoint constructs a Greeting endpoint.
func MakeGreetingEndpoint() server.JSONEndpoint {
	return func(r *http.Request) (int, interface{}, error) {
		// TODO - add real service call to get greeting
		greeting := "Gizmo Hello!!!"
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
