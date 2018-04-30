package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	proto "github.com/antklim/go-microservices/go-micro-greeter/pb"
	"github.com/micro/go-micro/client"
	web "github.com/micro/go-web"
)

func main() {
	service := web.NewService(
		web.Name("go-micro-web-greeter"),
	)

	service.HandleFunc("/greeting", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			var name string
			vars := r.URL.Query()
			names, exists := vars["name"]
			if !exists || len(names) != 1 {
				name = ""
			} else {
				name = names[0]
			}

			cl := proto.NewGreeterClient("go-micro-srv-greeter", client.DefaultClient)
			rsp, err := cl.Greeting(context.Background(), &proto.GreetingRequest{Name: name})
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			js, err := json.Marshal(rsp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	})

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
