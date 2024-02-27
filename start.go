package lambda

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

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
		lambdaContext := Context[Req]{
			status: http.StatusOK,
		}

		// For the string -> []byte we need to use a more effective way. For now, let's keep the naive approach.
		err := unmarshalBody(&gatewayReq, &lambdaContext)
		if err != nil {
			return c.errorHandler(err)
		}

		err = handler(ctx, &lambdaContext)
		if !errors.Is(err, &lambdaContext) && err != nil {
			lambdaContext.error = err
		}
		if lambdaContext.error != nil {
			return c.errorHandler(lambdaContext.error)
		}
		fmt.Println("lambdaContext.status", lambdaContext.status)
		fmt.Println("lambdaContext.header", lambdaContext.header)
		fmt.Println("lambdaContext.body.String()", lambdaContext.body.String())

		return events.APIGatewayProxyResponse{
			StatusCode:        lambdaContext.status,
			MultiValueHeaders: lambdaContext.header,
			Body:              lambdaContext.body.String(),
			IsBase64Encoded:   false,
		}, nil
	})
}

// unmarshalBody will unmarshal the body from the gateway request into the lambda context.
func unmarshalBody[Req any](gatewayReq *events.APIGatewayProxyRequest, lambdaContext *Context[Req]) error {
	fmt.Println("gatewayReq.HTTPMethod", gatewayReq.HTTPMethod)
	if gatewayReq.HTTPMethod == "GET" {
		return nil
	}
	fmt.Println("len(gatewayReq.Body)", len(gatewayReq.Body))
	if len(gatewayReq.Body) == 0 {
		return nil
	}
	var reader io.Reader
	fmt.Println("gatewayReq.IsBase64Encoded", gatewayReq.IsBase64Encoded)
	if gatewayReq.IsBase64Encoded {
		err := json.NewDecoder(
			base64.NewDecoder(
				base64.StdEncoding, bytes.NewReader([]byte(gatewayReq.Body)),
			),
		).Decode(&lambdaContext)
		if err != nil {
			return err
		}
	} else {
		reader = bytes.NewReader([]byte(gatewayReq.Body))
	}
	return json.NewDecoder(reader).Decode(lambdaContext.Request)
}
