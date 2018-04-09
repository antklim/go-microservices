package greeterservice

import (
	"context"

	"github.com/go-kit/kit/log"
)

// Service describe greetings service.
type Service interface {
	GetHealth() (bool, error)
	GetGreeting(ctx context.Context, name string) (string, error)
}

// New returns a basic Service with all of the expected middlewares wired in.
func New(logger log.Logger) Service {
	var svc Service
	{
		svc = NewGreeterService()
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// NewGreeterService returns implemetation of Service.
func NewGreeterService() Service {
	return greeterService{}
}

type greeterService struct{}

func (s greeterService) GetHealth() (bool, error) {
	return true, nil
}

func (s greeterService) GetGreeting(ctx context.Context, name string) (string, error) {
	greeting := "GO-KIT Hello " + name
	return greeting, nil
}
