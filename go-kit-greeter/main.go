package main

import (
	"context"
	// "errors"
	proto "github.com/antklim/go-microservices/go-kit-greeter/proto"
	"github.com/go-kit/kit/endpoint"
	// "strings"
)

type Greeter interface {
	Hello(ctx context.Context, in *proto.HelloRequest) (string, error)
}

type greeterService struct{}

func (g *greeterService) Hello(ctx context.Context, in *proto.HelloRequest) (string, error) {
	greeting := "GO-KIT Hello " + in.Name
	return greeting, nil
}

func makeHelloEndpoint(svc Greeter) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(proto.HelloRequest)
		greeting, err := svc.Hello(ctx, &req)
		if err != nil {
			return proto.HelloResponse{}, err
		}
		return proto.HelloResponse{Greeting: greeting}, nil
	}
}

func main() {
	svc := greeterService{}
	helloEndpoint := makeHelloEndpoint(&svc)
	helloEndpoint()
}
