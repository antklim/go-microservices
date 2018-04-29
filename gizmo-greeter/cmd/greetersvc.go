package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/antklim/go-microservices/gizmo-greeter/pkg/greeterendpoint"
	"github.com/antklim/go-microservices/gizmo-greeter/pkg/greetersd"
	"github.com/antklim/go-microservices/gizmo-greeter/pkg/greeterservice"
	"github.com/antklim/go-microservices/gizmo-greeter/pkg/greetertransport"
	"github.com/oklog/oklog/pkg/group"
)

func main() {
	var cfg *greetertransport.Config
	config.LoadJSONFile("./config.json", &cfg)
	server.SetConfigOverrides(cfg.Server)

	server.Init("gizmo-greeter", cfg.Server)

	var (
		service   = greeterservice.GreeterService{}
		endpoints = greeterendpoint.MakeServerEndpoints(service)
		registar  = greetersd.ConsulRegister("", "8500", "", strconv.Itoa(cfg.Server.HTTPPort))
	)

	var g group.Group
	{
		err := server.Register(greetertransport.NewTService(cfg, endpoints))
		if err != nil {
			server.Log.Fatal("unable to register service: ", err)
			os.Exit(1)
		}

		g.Add(func() error {
			registar.Register()
			return server.Run()
		}, func(err error) {
			registar.Deregister()
			server.Log.Fatal("server encountered a fatal error: ", err)
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	server.Log.Debug("exit", g.Run())
}
