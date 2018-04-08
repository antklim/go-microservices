package main

import (
	"flag"
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

	// var s Service {

	// }

	// var h http.Handler
	// {

	// }

	errs := make(chan error)

	logger.Log("exit", <-errs)
}
