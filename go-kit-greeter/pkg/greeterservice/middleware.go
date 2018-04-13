package greeterservice

import (
	"context"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service middleware.
type Middleware func(Service) Service

// LoggingMiddleware takes a logger as a dependency and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (m loggingMiddleware) Health() (healthy bool, err error) {
	defer func() {
		m.logger.Log("method", "Health", "healthy", healthy, "err", err)
	}()
	return m.next.Health()
}

func (m loggingMiddleware) Greeting(ctx context.Context, name string) (greeting string, err error) {
	defer func() {
		m.logger.Log("method", "Greeting", "name", name, "greeting", greeting, "err", err)
	}()
	return m.next.Greeting(ctx, name)
}
