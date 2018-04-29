package greetertransport

import (
	"context"

	"github.com/NYTimes/gizmo/server"
	"google.golang.org/grpc"

	"errors"
	"net/http"

	"github.com/NYTimes/gziphandler"
	pb "github.com/antklim/go-microservices/gizmo-greeter/pb"
	"github.com/antklim/go-microservices/gizmo-greeter/pkg/greeterendpoint"
	"github.com/sirupsen/logrus"
)

type (
	// TService will implement server.RPCService and handle all requests to the server.
	TService struct {
		Endpoints greeterendpoint.Endpoints
	}

	// Config is a struct to contain all the needed
	// configuration for our JSONService.
	Config struct {
		Server *server.Config
	}
)

// NewTService will instantiate a RPCService with the given configuration.
func NewTService(cfg *Config, endpoints greeterendpoint.Endpoints) *TService {
	return &TService{Endpoints: endpoints}
}

// Prefix returns the string prefix used for all endpoints within this service.
func (s *TService) Prefix() string {
	return ""
}

// Service provides the TService with a description of the service to serve and
// the implementation.
func (s *TService) Service() (*grpc.ServiceDesc, interface{}) {
	return &pb.Greeter_serviceDesc, s
}

// Middleware provides an http.Handler hook wrapped around all requests.
// In this implementation, we're using a GzipHandler middleware to
// compress our responses.
func (s *TService) Middleware(h http.Handler) http.Handler {
	return gziphandler.GzipHandler(h)
}

// ContextMiddleware provides a server.ContextHAndler hook wrapped around all
// requests. This could be handy if you need to decorate the request context.
func (s *TService) ContextMiddleware(h server.ContextHandler) server.ContextHandler {
	return h
}

// JSONMiddleware provides a JSONEndpoint hook wrapped around all requests.
// In this implementation, we're using it to provide application logging and to check errors
// and provide generic responses.
func (s *TService) JSONMiddleware(j server.JSONContextEndpoint) server.JSONContextEndpoint {
	return func(ctx context.Context, r *http.Request) (int, interface{}, error) {

		status, res, err := j(ctx, r)
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

// ContextEndpoints may be needed if your server has any non-RPC-able
// endpoints. In this case, we have none but still need this method to
// satisfy the server.RPCService interface.
func (s *TService) ContextEndpoints() map[string]map[string]server.ContextHandlerFunc {
	return map[string]map[string]server.ContextHandlerFunc{}
}

// JSONEndpoints is a listing of all endpoints available in the TService.
func (s *TService) JSONEndpoints() map[string]map[string]server.JSONContextEndpoint {
	return map[string]map[string]server.JSONContextEndpoint{
		"/health": map[string]server.JSONContextEndpoint{
			"GET": s.Endpoints.HealthEndpoint,
		},
		"/greeting": map[string]server.JSONContextEndpoint{
			"GET": s.Endpoints.GreetingEndpoint,
		},
	}
}
