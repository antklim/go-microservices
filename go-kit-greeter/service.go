package gokitgreeter

import (
	"context"
)

type Service interface {
	GetHealth(ctx context.Context) (bool, error)
	GetGreeting(ctx context.Context, name string) (string, error)
}

type greeterService struct{}

func (s *greeterService) GetHealth(ctx context.Context) (bool, error) {
	return true, nil
}

func (s *greeterService) GetGreeting(ctx context.Context, name string) (string, error) {
	greeting := "GO-KIT Hello " + name
	return greeting, nil
}
