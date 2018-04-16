package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"

	"github.com/antklim/go-microservices/go-kit-greeter/pkg/greetersd"

	"github.com/antklim/go-microservices/go-kit-greeter/pkg/greeterendpoint"
	"github.com/antklim/go-microservices/go-kit-greeter/pkg/greeterservice"
	"github.com/antklim/go-microservices/go-kit-greeter/pkg/greetertransport"

	"github.com/go-kit/kit/log"
	"github.com/oklog/oklog/pkg/group"
)

func main() {
	fs := flag.NewFlagSet("greetersvc", flag.ExitOnError)
	var (
		debugAddr  = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
		consulAddr = fs.String("consul.addr", "", "Consul Address")
		consulPort = fs.String("consul.port", "8500", "Consul Port")
		httpAddr   = fs.String("http.addr", "", "HTTP Listen Address")
		httpPort   = fs.String("http.port", "8081", "HTTP Listen Port")
		// grpcAddr  = fs.String("grpc-addr", ":8082", "gRPC listen address")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var service greeterservice.Service
	{
		service = greeterservice.GreeterService{}
		service = greeterservice.LoggingMiddleware(logger)(service)
	}

	var (
		endpoints   = greeterendpoint.MakeServerEndpoints(service, logger)
		httpHandler = greetertransport.NewHTTPHandler(endpoints, logger)
		registar    = greetersd.ConsulRegister(*consulAddr, *consulPort, *httpAddr, *httpPort)
		// grpcServer  = greetertransport.NewGRPCServer(endpoints, logger)
	)

	var g group.Group
	{
		// The debug listener mounts the http.DefaultServeMux, and serves up
		// stuff like the Go debug and profiling routes, and so on.
		debugListener, err := net.Listen("tcp", *debugAddr)
		if err != nil {
			logger.Log("transport", "debug/HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "debug/HTTP", "addr", *debugAddr)
			return http.Serve(debugListener, http.DefaultServeMux)
		}, func(error) {
			debugListener.Close()
		})
	}
	{
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", *httpAddr, "port", *httpPort)
			registar.Register()
			handler := httpHandler
			return http.ListenAndServe(":"+*httpPort, handler)
		}, func(error) {
			registar.Deregister()
		})
	}
	// {
	// 	// The HTTP listener mounts the Go kit HTTP handler we created.
	// 	httpListener, err := net.Listen("tcp", *httpAddr)
	// 	if err != nil {
	// 		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
	// 		os.Exit(1)
	// 	}
	// 	g.Add(func() error {
	// 		logger.Log("transport", "HTTP", "addr", *httpAddr)
	// 		return http.Serve(httpListener, httpHandler)
	// 	}, func(error) {
	// 		httpListener.Close()
	// 	})
	// }
	// {
	// 	// The gRPC listener mounts the Go kit gRPC server we created.
	// 	grpcListener, err := net.Listen("tcp", *grpcAddr)
	// 	if err != nil {
	// 		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
	// 		os.Exit(1)
	// 	}
	// 	g.Add(func() error {
	// 		logger.Log("transport", "gRPC", "addr", *grpcAddr)
	// 		baseServer := grpc.NewServer()
	// 		pb.RegisterGreeterServer(baseServer, grpcServer)
	// 		return baseServer.Serve(grpcListener)
	// 	}, func(error) {
	// 		grpcListener.Close()
	// 	})
	// }
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
	logger.Log("exit", g.Run())
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}
