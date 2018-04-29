package main

import (
	"context"
	"fmt"

	proto "github.com/antklim/go-microservices/go-micro-greeter/pb"
	micro "github.com/micro/go-micro"
)

func main() {
	service := micro.NewService(micro.Name("go-micro-srv-greeter.client"))
	service.Init()

	client := proto.NewGreeterClient("go-micro-srv-greeter", service.Client())

	rsp, err := client.Greeting(context.Background(), &proto.GreetingRequest{Name: "gomicro RPC call"})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(rsp.Greeting)
}
