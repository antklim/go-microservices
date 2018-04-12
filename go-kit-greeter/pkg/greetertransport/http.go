package greetertransport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/antklim/go-microservices/go-kit-greeter/pkg/greeterendpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// NewHTTPHandler returns an HTTP handler that makes a set of endpoints
// available on predefined paths.
func NewHTTPHandler(endpoints greeterendpoint.Endpoints, logger log.Logger) http.Handler {
	m := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerErrorLogger(logger),
	}

	// GET /health      retrieves service heath information
	// GET /hello?name  retrieves greeting

	m.Methods("GET").Path("/health").Handler(httptransport.NewServer(
		endpoints.GetHealthEndpoint,
		decodeHTTPGetHealthRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	m.Methods("GET").Path("/hello").Handler(httptransport.NewServer(
		endpoints.GetGreetingEndpoint,
		decodeHTTPGetGreeterRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	return m
}

func decodeHTTPGetHealthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req greeterendpoint.GetHealthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPGetGreeterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req greeterendpoint.GetGreetingRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func err2code(err error) int {
	switch err {
	default:
		return http.StatusInternalServerError
	}
}

type errorWrapper struct {
	Error string `json:"error"`
}

// encodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(greeterendpoint.Failer); ok && f.Failed() != nil {
		encodeError(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
