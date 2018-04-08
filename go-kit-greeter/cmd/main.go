package main

import (
	"flag"
	"fmt"
	gokitgreeter "github.com/antklim/go-microservices/go-kit-greeter"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	var s gokitgreeter.Service
	{
		s = gokitgreeter.NewGoKitGreeterService()
	}

	var h http.Handler
	{
		h = gokitgreeter.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("trasport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
