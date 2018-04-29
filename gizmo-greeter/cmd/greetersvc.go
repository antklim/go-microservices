package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"text/tabwriter"

	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/antklim/go-microservices/gizmo-greeter/pkg/greeterendpoint"
	"github.com/antklim/go-microservices/gizmo-greeter/pkg/greetersd"
	"github.com/antklim/go-microservices/gizmo-greeter/pkg/greeterservice"
	"github.com/antklim/go-microservices/gizmo-greeter/pkg/greetertransport"
	"github.com/oklog/oklog/pkg/group"
)

func main() {
	fs := flag.NewFlagSet("greetersvc", flag.ExitOnError)
	var (
		configPath = fs.String("config.path", "./config.json", "Config file path")
		consulAddr = fs.String("consul.addr", "", "Consul Address")
		consulPort = fs.String("consul.port", "8500", "Consul Port")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])

	var cfg *greetertransport.Config
	config.LoadJSONFile(*configPath, &cfg)

	server.Init("gizmo-greeter", cfg.Server)

	var (
		service   = greeterservice.GreeterService{}
		endpoints = greeterendpoint.MakeServerEndpoints(service)
		registar  = greetersd.ConsulRegister(*consulAddr, *consulPort, "", strconv.Itoa(cfg.Server.HTTPPort))
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
