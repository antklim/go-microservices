#!/usr/bin/env sh

echo "Stopping go-micro-greeter"
kill `cat go-micro-greeter.pid`
sleep 3

echo "Stopping go-kit-greeter"
kill `cat go-kit-greeter.pid`
sleep 3

echo "Stopping Consul"
kill `cat consulagent.pid`
