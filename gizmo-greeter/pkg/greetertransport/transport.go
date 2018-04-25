package greetertransport

import (
	// "github.com/NYTimes/gizmo/server"

	"errors"
	"net/http"

	"../../../../gizmo/server"
	"github.com/NYTimes/gziphandler"
	"github.com/sirupsen/logrus"
)

type (
	// JSONService will implement server.JSONService and
	// handle all requests to the server.
	JSONService struct{}

	// Config is a struct to contain all the needed
	// configuration for our JSONService.
	Config struct {
		Server *server.Config
	}
)

////////////////

// Prefix returns the string prefix used for all endpoints within this service.
func (s *JSONService) Prefix() string {
	return ""
}

// NewJSONService will instantiate a JSONService with the given configuration.
func NewJSONService(cfg *Config) *JSONService {
	return &JSONService{}
}

// Middleware provides an http.Handler hook wrapped around all requests.
// In this implementation, we're using a GzipHandler middleware to
// compress our responses.
func (s *JSONService) Middleware(h http.Handler) http.Handler {
	return gziphandler.GzipHandler(h)
}

// JSONMiddleware provides a JSONEndpoint hook wrapped around all requests.
// In this implementation, we're using it to provide application logging and to check errors
// and provide generic responses.
func (s *JSONService) JSONMiddleware(j server.JSONEndpoint) server.JSONEndpoint {
	return func(r *http.Request) (int, interface{}, error) {

		status, res, err := j(r)
		if err != nil {
			server.LogWithFields(r).WithFields(logrus.Fields{
				"error": err,
			}).Error("problems with serving request")
			return http.StatusServiceUnavailable, nil, errors.New("sorry, this service is unavailable")
		}

		server.LogWithFields(r).Info("success!")
		return status, res, nil
	}
}

// JSONEndpoints is a listing of all endpoints available in the JSONService.
func (s *JSONService) JSONEndpoints() map[string]map[string]server.JSONEndpoint {
	return map[string]map[string]server.JSONEndpoint{
		"/health": map[string]server.JSONEndpoint{
			"GET": s.Health,
		},
		"/greeting": map[string]server.JSONEndpoint{
			"GET": s.Greeting,
		},
	}
}

// Health endpoint.
func (s *JSONService) Health(r *http.Request) (int, interface{}, error) {
	// healthy := service.Health() should call real service instance
	return http.StatusOK, HealthResponse{Healthy: true}, nil
}

// Greeting endpoint.
func (s *JSONService) Greeting(r *http.Request) (int, interface{}, error) {
	greeting := "Gizmo Hello!!!"
	return http.StatusOK, GreetingResponse{Greeting: greeting}, nil
}

//////////////

// HealthRequest collects the request parameters for the Health method.
type HealthRequest struct{}

// HealthResponse collects the response values for the Health method.
type HealthResponse struct {
	Healthy bool  `json:"healthy,omitempty"`
	Err     error `json:"err,omitempty"`
}

// GreetingRequest collects the request parameters for the Greeting method.
type GreetingRequest struct {
	Name string `json:"name,omitempty"`
}

// GreetingResponse collects the response values for the Greeting method.
type GreetingResponse struct {
	Greeting string `json:"greeting,omitempty"`
	Err      error  `json:"err,omitempty"`
}
