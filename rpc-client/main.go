package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	gizmoservice "github.com/antklim/go-microservices/gizmo-greeter/pb"
	gokitservice "github.com/antklim/go-microservices/go-kit-greeter/pb"
	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
)

func main() {
	fs := flag.NewFlagSet("greeterclient", flag.ExitOnError)
	var (
		// gomicroAddr = fs.String("gomicro.addr", "", "The Go Micro greeter service address")
		gokitAddr = fs.String("gokit.addr", "127.0.0.1:9120", "The Go Kit greeter service address")
		gizmoAddr = fs.String("gizmo.port", "127.0.0.1:9220", "The Gizmo greeter service address")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// go-kit client =============================================================
	gokitConn, err := grpc.Dial(*gokitAddr, grpc.WithInsecure())
	if err != nil {
		logger.Log("grpcGoKitConnectionErr", err)
		os.Exit(1)
	}
	defer func() {
		err := gokitConn.Close()
		if err != nil {
			logger.Log("goKitConnectionError", err)
		}
	}()

	gokitClient := gokitservice.NewGreeterClient(gokitConn)
	gokitServiceResponse, err := gokitClient.Greeting(context.Background(), &gokitservice.GreetingRequest{Name: "gokit RPC call"})
	if err != nil {
		logger.Log("goKitServiceErr", err)
		return
	}
	logger.Log("goKitResponse", gokitServiceResponse.Greeting)

	// gizmo client =============================================================
	gizmoConn, err := grpc.Dial(*gizmoAddr, grpc.WithInsecure())
	if err != nil {
		logger.Log("grpcGizmoConnectionErr", err)
		os.Exit(1)
	}
	defer func() {
		err := gizmoConn.Close()
		if err != nil {
			logger.Log("gizmoConnectionError", err)
		}
	}()

	gizmoClient := gizmoservice.NewGreeterClient(gizmoConn)
	gizmoServiceResponse, err := gizmoClient.Greeting(context.Background(), &gizmoservice.GreetingRequest{Name: "gizmo RPC call"})
	if err != nil {
		logger.Log("gizmoServiceErr", err)
		return
	}
	logger.Log("gizmoResponse", gizmoServiceResponse.Greeting)
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
