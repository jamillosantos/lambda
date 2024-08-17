package http

import (
	"context"
	"encoding/json"
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
	errorHandler func(error) (HttpResponse, error)
}

func defaultOpts() options {
	return options{
		resources:    make([]Resource, 0),
		errorHandler: DefaultHttpErrorHandler,
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
func WithErrorHandler(h func(error) (HttpResponse, error)) HttpOption {
	return func(o *options) {
		o.errorHandler = h
	}
}

var DefaultHttpErrorHandler = func(err error) (HttpResponse, error) {
	e, mErr := json.Marshal(err.Error())
	if mErr != nil {
		return HttpResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       []byte(`"failed to marshal error"`),
		}, mErr

	}
	return HttpResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       e,
	}, nil
}
