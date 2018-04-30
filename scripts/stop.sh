#!/usr/bin/env sh

echo "Stopping go-micro-greeter-srv"
kill -9 `cat go-micro-greeter-srv.pid`
sleep 3

echo "Stopping go-micro-greeter-web"
kill -9 `cat go-micro-greeter-web.pid`
sleep 3

echo "Stopping go-kit-greeter"
kill -9 `cat go-kit-greeter.pid`
sleep 3

echo "Stopping gizmo-greeter"
kill -9 `cat gizmo-greeter.pid`
sleep 3

echo "Stopping Consul"
kill -9 `cat consulagent.pid`
