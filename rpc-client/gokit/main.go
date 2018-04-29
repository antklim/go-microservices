package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	gokitservice "github.com/antklim/go-microservices/go-kit-greeter/pb"
	"google.golang.org/grpc"
)

func main() {
	fs := flag.NewFlagSet("greeterclient", flag.ExitOnError)
	var (
		serviceAddr = fs.String("sercice.addr", "127.0.0.1:9120", "The Go Kit greeter service address")
		name        = fs.String("name", "gokit RPC call", "The Name to greet")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])

	conn, err := grpc.Dial(*serviceAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("grpcGoKitConnectionErr", err)
		os.Exit(1)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("goKitConnectionError", err)
		}
	}()

	client := gokitservice.NewGreeterClient(conn)
	serviceResponse, err := client.Greeting(context.Background(), &gokitservice.GreetingRequest{Name: *name})
	if err != nil {
		fmt.Println("goKitServiceErr", err)
		return
	}
	fmt.Println("goKitResponse", serviceResponse.Greeting)
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
