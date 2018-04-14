package greetertransport

import (
	"context"

	"github.com/antklim/go-microservices/go-kit-greeter/pb"
	"github.com/antklim/go-microservices/go-kit-greeter/pkg/greeterendpoint"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	oldcontext "golang.org/x/net/context"
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
			endpoints.GreetingEndpoint,
			decodeGRPCGreetingRequest,
			encodeGRPCGreetingResponse,
			options...,
		),
	}
}

// Greeting implementation of the method of the GreeterService interface.
func (s *grpcServer) Greeting(ctx oldcontext.Context, req *pb.GreetingRequest) (*pb.GreetingResponse, error) {
	_, res, err := s.greeter.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GreetingResponse), nil
}

// decodeGRPCGreetingRequest is a transport/grpc.DecodeRequestFunc that converts
// a gRPC greeting request to a user-domain greeting request.
func decodeGRPCGreetingRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GreetingRequest)
	return greeterendpoint.GreetingRequest{Name: req.Name}, nil
}

// encodeGRPCGreetingResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain greeting response to a gRPC greeting response.
func encodeGRPCGreetingResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(greeterendpoint.GreetingResponse)
	return &pb.GreetingResponse{Greeting: res.Greeting}, nil
}
