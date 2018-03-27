package main

import (
	proto "github.com/antklim/go-microservices/go-micro-greeter/proto"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"log"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, in *proto.HelloRequest, out *proto.HelloResponse) error {
	out.Greeting = "GO-MICRO Hello " + in.Name
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.greeter"),
		micro.Version("latest"),
	)

	service.Init()

	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
