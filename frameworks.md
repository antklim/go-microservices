https://github.com/nytimes/gizmo
https://github.com/koding/kite
https://github.com/go-kit/kit
https://github.com/micro/go-micro
---
https://www.consul.io/
https://www.consul.io/docs/guides/consul-containers.html
https://www.hashicorp.com/blog/official-consul-docker-image


# Microservice in go-micro
## Defining API interface
1. Get protobuf generator
```
$ go get github.com/micro/protoc-gen-micro
```

2. Define the service as greeter.proto
3. Generate interface
```
$ protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. greeter.proto
```
