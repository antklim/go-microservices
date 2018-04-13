package greetertransport

import (
	"github.com/antklim/go-microservices/go-kit-greeter/pb"
	"github.com/antklim/go-microservices/go-kit-greeter/pkg/greeterendpoint"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	greeter grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC GreeterServer.
func NewGRPCServer(endpoints greeterendpoint.Endpoints, logger log.Logger) pb.GreeterServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}

	return &grpcServer{
		greeter: grpctransport.NewServer(
			endpoints.GetGreetingEndpoint,
			decodeGRPCGetGreetingRequest,
			encodeGRPCGetGreetingResponse,
			options...,
		),
	}
}

func (s *grpcServer) Hello()
