package lambda

import (
	"context"
	"encoding/json"
	"net/http"
)

// Resource is an interface that represents a resource that can be started and stopped.
// Example: A connection to the database or message broker, etc;
type Resource interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type options struct {
	resources    []Resource
	errorHandler func(error) (APIGatewayProxyResponse, error)
}

func defaultOpts() options {
	return options{
		resources:    make([]Resource, 0),
		errorHandler: DefaultErrorHandler,
	}
}

type Option func(*options)

// WithResources is an option that allows you to pass resources to the lambda function.
func WithResources(r ...Resource) Option {
	return func(o *options) {
		o.resources = append(o.resources, r...)
	}
}

// WithErrorHandler is an option that allows you to pass a custom error handler to the lambda function.
func WithErrorHandler(h func(error) (APIGatewayProxyResponse, error)) Option {
	return func(o *options) {
		o.errorHandler = h
	}
}

var DefaultErrorHandler = func(err error) (APIGatewayProxyResponse, error) {
	e, mErr := json.Marshal(err.Error())
	if mErr != nil {
		return APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       []byte(`"failed to marshal error"`),
		}, mErr

	}
	return APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       e,
	}, nil
}
