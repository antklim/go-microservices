package greeterservice

import (
	"time"

	"github.com/go-kit/kit/log"
)

// ServiceMiddleware describes a service middleware.
type ServiceMiddleware func(Service) Service

// LoggingMiddleware takes a logger as a dependency and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return loggingMiddleware{next, logger}
	}
}

type loggingMiddleware struct {
	Service
	logger log.Logger
}

func (m loggingMiddleware) Health() (healthy bool) {
	defer func(begin time.Time) {
		m.logger.Log(
			"method", "Health",
			"healthy", healthy,
			"took", time.Since(begin),
		)
	}(time.Now())
	healthy = m.Service.Health()
	return
}

func (m loggingMiddleware) Greeting(name string) (greeting string) {
	defer func(begin time.Time) {
		m.logger.Log(
			"method", "Greeting",
			"name", name,
			"greeting", greeting,
			"took", time.Since(begin),
		)
	}(time.Now())
	greeting = m.Service.Greeting(name)
	return
}
