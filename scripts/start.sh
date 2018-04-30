#!/usr/bin/env sh

echo "Starting Consul"
nohup consul agent -dev -ui &
echo $! > consulagent.pid
sleep 5

echo "Starting go-micro-greeter srv"
nohup go run ../go-micro-greeter/srv/main.go &
echo $! > go-micro-greeter-srv.pid
sleep 1

echo "Starting go-micro-greeter web"
nohup go run ../go-micro-greeter/web/main.go &
echo $! > go-micro-greeter-web.pid
sleep 1

echo "Starting go-kit-greeter"
nohup go run ../go-kit-greeter/cmd/greetersvc.go &
echo $! > go-kit-greeter.pid
sleep 1

echo "Starting gizmo-greeter"
nohup go run ../gizmo-greeter/cmd/greetersvc.go -config.path='/Users/antklim/code/HOME/go/go-microservices/gizmo-greeter/cmd/config.json' &
echo $! > gizmo-greeter.pid
sleep 1
