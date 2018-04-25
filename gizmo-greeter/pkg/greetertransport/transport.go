package greetertransport

import (
	"github.com/NYTimes/gizmo/server"

	"errors"
	"net/http"

	endpoints "../greeterendpoint"
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

var ep = endpoints.MakeServerEndpoints()

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
			"GET": ep.HealthEndpoint,
		},
		"/greeting": map[string]server.JSONEndpoint{
			"GET": ep.GreetingEndpoint,
		},
	}
}
