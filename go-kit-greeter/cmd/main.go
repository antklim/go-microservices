package main

import (
	"flag"
	gokitgreeter "github.com/antklim/go-microservices/go-kit-greeter"
	"github.com/go-kit/kit/log"
	"os"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP Listen Address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// var s gokitgreeter.Service
	// {
	// 	s = gokitgreeter.NewGoKitGreeterService()
	// }

	var h http.Handler
	{
		h = gokitgreeter.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)

	logger.Log("exit", <-errs)
}
