package lambda

type options[Resp any] struct {
	resources    []Resource
	errorHandler func(error) (Resp, error)
}

func defaultOpts[Resp any]() options[Resp] {
	return options[Resp]{
		resources:    make([]Resource, 0),
		errorHandler: nil,
	}
}

type Option[Resp any] func(*options[Resp])

// WithResources is an option that allows you to pass resources to the lambda function.
func WithResources[Resp any](r ...Resource) Option[Resp] {
	return func(o *options[Resp]) {
		o.resources = append(o.resources, r...)
	}
}

// WithErrorHandler is an option that allows you to pass a custom error handler to the lambda function.
func WithErrorHandler[Resp any](h func(error) (Resp, error)) Option[Resp] {
	return func(o *options[Resp]) {
		o.errorHandler = h
	}
}
