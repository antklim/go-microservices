package greetertransport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/antklim/go-microservices/go-kit-greeter/pkg/greeterendpoint"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler")
)

// NewHTTPHandler returns an HTTP handler that makes a set of endpoints
// available on predefined paths.
func NewHTTPHandler(endpoints greeterendpoint.Endpoints, logger log.Logger) http.Handler {
	m := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerErrorLogger(logger),
	}

	// GET /health         retrieves service heath information
	// GET /greeting?name  retrieves greeting

	m.Methods("GET").Path("/health").Handler(httptransport.NewServer(
		endpoints.HealthEndpoint,
		decodeHTTPHealthRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	m.Methods("GET").Path("/greeting").Handler(httptransport.NewServer(
		endpoints.GreetingEndpoint,
		decodeHTTPGreetingRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	return m
}

// MakeHTTPClientEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the remote instance, via a transport/http.Client.
func MakeHTTPClientEndpoints(instance string) (greeterendpoint.Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return greeterendpoint.Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	var healthEndpoint endpoint.Endpoint
	{
		healthEndpoint = httptransport.NewClient(
			"GET",
			tgt,
			encodeHTTPHealthRequest,
			decodeHTTPHealthResponse,
			options...).Endpoint()
	}

	var greetingEndpoint endpoint.Endpoint
	{
		greetingEndpoint = httptransport.NewClient(
			"GET",
			tgt,
			encodeHTTPGreetingRequest,
			decodeHTTPGreetingResponse,
			options...).Endpoint()
	}

	return greeterendpoint.Endpoints{
		HealthEndpoint:   healthEndpoint,
		GreetingEndpoint: greetingEndpoint,
	}, nil
}

func encodeHTTPHealthRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").path("/health")
	req.Method, req.URL.Path = "GET", "/health"
	return encodeHTTPGenericRequest(ctx, req, request)
}

func decodeHTTPHealthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := greeterendpoint.HealthRequest{}
	return req, nil
}

func decodeHTTPHealthResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response greeterendpoint.HealthResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func encodeHTTPGreetingRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").path("/greeting?name=bob")
	req.Method, req.URL.Path, req.URL.RawQuery = "GET", "/greeting", "?name=bob"
	return encodeHTTPGenericRequest(ctx, req, request)
}

func decodeHTTPGreetingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := r.URL.Query()
	names, exists := vars["name"]
	if !exists || len(names) != 1 {
		return nil, ErrBadRouting
	}
	req := greeterendpoint.GreetingRequest{Name: names[0]}
	return req, nil
}

func decodeHTTPGreetingResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response greeterendpoint.GreetingResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// encodeHTTPGenericRequest likewise JSON-encodes the request to the HTTP request body.
func encodeHTTPGenericRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
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
