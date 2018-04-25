package main

import (
	server "../../../gizmo/server"
	service "../pkg/greetertransport"
	"github.com/NYTimes/gizmo/config"
)

func main() {
	var cfg *service.Config
	config.LoadJSONFile("./config.json", &cfg)

	server.Init("gizmo-hello-world", cfg.Server)

	err := server.Register(service.NewJSONService(cfg))
	if err != nil {
		server.Log.Fatal("unable to register service: ", err)
	}

	err = server.Run()
	if err != nil {
		server.Log.Fatal("server encountered a fatal error: ", err)
	}
}
