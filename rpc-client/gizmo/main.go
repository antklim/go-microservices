package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	gizmoservice "github.com/antklim/go-microservices/gizmo-greeter/pb"
	"google.golang.org/grpc"
)

func main() {
	fs := flag.NewFlagSet("greeterclient", flag.ExitOnError)
	var (
		serviceAddr = fs.String("service.addr", "127.0.0.1:9220", "The Gizmo greeter service address")
		name        = fs.String("name", "gizmo RPC call", "The Name to greet")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])

	conn, err := grpc.Dial(*serviceAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("grpcGizmoConnectionErr", err)
		os.Exit(1)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("gizmoConnectionError", err)
		}
	}()

	client := gizmoservice.NewGreeterClient(conn)
	serviceResponse, err := client.Greeting(context.Background(), &gizmoservice.GreetingRequest{Name: *name})
	if err != nil {
		fmt.Println("gizmoServiceErr", err)
		return
	}
	fmt.Println("gizmoResponse", serviceResponse.Greeting)
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
