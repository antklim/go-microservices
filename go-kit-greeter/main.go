package main

import (
	"context"
	"errors"
	"strings"
)

type Greeter interface {
	Hello(context.Context, string) (string, error)
}

type greeterService struct{}

func (greeterService) Hello(_ context.Context, s string) (string, error) {
	return "GO-KIT Hello " + s
}
