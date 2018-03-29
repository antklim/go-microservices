package gokitgreeter

import (
	"errors"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func decodeGetHealthRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return getHealthRequest{}, nil
}

func decodeGetGreeterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, ErrBadRouting
	}

	return getGreetingRequest{Name: name}, nil
}
