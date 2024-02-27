package lambda

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Resource is an interface that represents a resource that can be started and stopped.
// Example: A connection to the database or message broker, etc;
type Resource interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type options struct {
	resources    []Resource
	errorHandler func(error) (events.APIGatewayProxyResponse, error)
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
func WithErrorHandler(h func(error) (events.APIGatewayProxyResponse, error)) Option {
	return func(o *options) {
		o.errorHandler = h
	}
}

var DefaultErrorHandler = func(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       err.Error(),
	}, nil
}
