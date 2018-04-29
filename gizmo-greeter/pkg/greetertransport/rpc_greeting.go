package greetertransport

import (
	pb "github.com/antklim/go-microservices/gizmo-greeter/pb"
	ocontext "golang.org/x/net/context"
)

// Greeting implementation of the gRPC service.
func (s *TService) Greeting(ctx ocontext.Context, r *pb.GreetingRequest) (*pb.GreetingResponse, error) {
	return &pb.GreetingResponse{Greeting: "Hola Gizmo RPC " + r.Name}, nil
}
