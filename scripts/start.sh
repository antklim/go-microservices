#!/usr/bin/env sh

nohup consul agent -dev -ui &
echo $! > consulagent.pid
sleep 5

nohup go run ../go-micro-greeter/main.go &
echo $! > go-micro-greeter.pid
sleep 1

nohup go run ../go-kit-greeter/cmd/greetersvc.go &
echo $! > go-kit-greeter.pid
sleep 1

nohup go run ../gizmo-greeter/cmd/greetersvc.go -config.path='/Users/antklim/code/HOME/go/go-microservices/gizmo-greeter/cmd/config.json' &
echo $! > gizmo-greeter.pid
sleep 1
