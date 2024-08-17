package lambda

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

type Handler[Req any, Resp any] func(*Context[Req]) (Resp, error)

// Start will start the lambda function with the given handler and httpOptions.
//
// Start should be used when not using http events.
func Start[Req any, Resp any](handler Handler[Req, Resp], opts ...Option[Resp]) {
	c := defaultOpts[Resp]()
	for _, o := range opts {
		o(&c)
	}

	startCtx := context.Background()
	for _, r := range c.resources {
		if err := r.Start(startCtx); err != nil {
			panic(fmt.Errorf("failed to start resource %s: %w", r.Name(), err))
		}
	}

	lambda.Start(func(ctx context.Context, request Req) (Resp, error) {
		lambdaContext := Context[Req]{
			Context: ctx,
			Request: request,
			Locals:  make(map[string]any),
		}

		resp, err := handler(&lambdaContext)
		if err != nil {
			return resp, err
		}

		return resp, nil
	})
}
