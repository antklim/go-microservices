#!/usr/bin/env sh

echo "Stopping go-micro-greeter"
kill -9 `cat go-micro-greeter.pid`
sleep 3

echo "Stopping go-kit-greeter"
kill -9 `cat go-kit-greeter.pid`
sleep 3

echo "Stopping gizmo-greeter"
kill -9 `cat gizmo-greeter.pid`
sleep 3

echo "Stopping Consul"
kill -9 `cat consulagent.pid`
