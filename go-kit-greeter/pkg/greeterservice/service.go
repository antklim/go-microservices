package greeterservice

// Service describe greetings service.
type Service interface {
	Health() bool
	Greeting(name string) string
}

// GreeterService implementation of the Service interface.
type GreeterService struct{}

// Health implementation of the Service.
func (GreeterService) Health() bool {
	return true
}

// Greeting implementation of the Service.
func (GreeterService) Greeting(name string) (greeting string) {
	greeting = "GO-KIT Hello " + name
	return
}
