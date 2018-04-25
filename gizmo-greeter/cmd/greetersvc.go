package main

import (
	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/antklim/go-microservices/gizmo-greeter/pkg/greeterendpoint"
	"github.com/antklim/go-microservices/gizmo-greeter/pkg/greeterservice"
	"github.com/antklim/go-microservices/gizmo-greeter/pkg/greetertransport"
)

func main() {
	var cfg *greetertransport.Config
	config.LoadJSONFile("./config.json", &cfg)

	server.Init("gizmo-greeter", cfg.Server)

	var service greeterservice.Service
	{
		service = greeterservice.GreeterService{}
	}

	var endpoints = greeterendpoint.MakeServerEndpoints(service)

	err := server.Register(greetertransport.NewJSONService(cfg, endpoints))
	if err != nil {
		server.Log.Fatal("unable to register service: ", err)
	}

	err = server.Run()
	if err != nil {
		server.Log.Fatal("server encountered a fatal error: ", err)
	}
}
