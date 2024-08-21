package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type Resource interface {
	Name() string
	Start(context.Context) error
}

type HttpResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       json.RawMessage
}

type options struct {
	resources    []Resource
	errorHandler func(context.Context, error) (HttpResponse, error)
}

func defaultOpts() options {
	return options{
		resources:    make([]Resource, 0),
		errorHandler: DefaultErrorHandler,
	}
}

type HttpOption func(*options)

// WithResources is an option that allows you to pass resources to the lambda function.
func WithResources(r ...Resource) HttpOption {
	return func(o *options) {
		o.resources = append(o.resources, r...)
	}
}

// WithErrorHandler is an option that allows you to pass a custom error handler to the lambda function.
func WithErrorHandler(h func(context.Context, error) (HttpResponse, error)) HttpOption {
	return func(o *options) {
		o.errorHandler = h
	}
}

// Error is a struct that implements ErrorResponse. It represents an error that can be returned by the lambda function.
type Error struct {
	StatusCode int
	Headers    map[string]string
	Message    string
}

func (h *Error) Error() string {
	return h.Message
}

func (h *Error) HttpStatusCode() int {
	return h.StatusCode
}

func (h *Error) HttpHeaders() map[string]string {
	return h.Headers
}

type httpErrorBody struct {
	Message string `json:"message,omitempty"`
}

func (h *Error) HttpBody() (json.RawMessage, error) {
	b, err := json.Marshal(httpErrorBody{
		Message: h.Message,
	})
	if err != nil {
		return nil, err
	}
	return b, nil
}

type ErrorResponse interface {
	HttpStatusCode() int
	HttpHeaders() map[string]string
	HttpBody() (json.RawMessage, error)
}

// DefaultErrorHandler is the default error handler for the lambda function. If th error implements ErrorResponse, it
// will return an HttpResponse with the status code, headers and body. If the error does not implement ErrorResponse, it
// will return a generic internal server error message.
//
// The original error message is not returned to avoid accidentally leaking sensitive information.
var DefaultErrorHandler = func(_ context.Context, err error) (HttpResponse, error) {
	e := err
	for {
		if httpErr, ok := e.(ErrorResponse); ok {
			b, err := httpErr.HttpBody()
			if err != nil {
				break // Exits the loop and returns the original error.
			}
			return HttpResponse{
				StatusCode: httpErr.HttpStatusCode(),
				Headers:    httpErr.HttpHeaders(),
				Body:       b,
			}, nil
		}
		if e = errors.Unwrap(e); e == nil {
			break
		}
	}
	return HttpResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       []byte(`{"message":"Internal Server Error"}`),
	}, nil
}
