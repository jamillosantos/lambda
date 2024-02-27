package lambda

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Start will start the lambda function with the given handler and options.
//
// The handler is a function that receives a context and a pointer to a Context.
func Start[Req any](handler func(context.Context, *Context[Req]) error, opts ...Option) {
	c := defaultOpts()
	for _, o := range opts {
		o(&c)
	}

	lambda.Start(func(ctx context.Context, gatewayReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var lambdaContext Context[Req]

		// For the string -> []byte we need to use a more effective way. For now, let's keep the naive approach.
		if gatewayReq.IsBase64Encoded {
			err := json.NewDecoder(
				base64.NewDecoder(
					base64.StdEncoding, bytes.NewReader([]byte(gatewayReq.Body)),
				),
			).Decode(&lambdaContext)
			if err != nil {
				return c.errorHandler(err)
			}
		} else {
			err := json.Unmarshal([]byte(gatewayReq.Body), &lambdaContext)
			if err != nil {
				return c.errorHandler(err)
			}
		}

		res := handler(ctx, &lambdaContext)
		if lambdaContext.error != nil {
			return c.errorHandler(lambdaContext.error)
		}

		d, err := json.Marshal(res)
		if err != nil {
			return c.errorHandler(err)
		}

		return events.APIGatewayProxyResponse{
			StatusCode:        lambdaContext.status,
			MultiValueHeaders: lambdaContext.header,
			Body:              string(d),
			IsBase64Encoded:   false,
		}, nil
	})
}
