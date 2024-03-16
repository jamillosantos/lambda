package lambda

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
)

type Handler[Req any, Resp any] func(*Context[Req, Resp]) error

// Start will start the lambda function with the given handler and options.
//
// The handler is a function that receives a context and a pointer to a Context.
func Start[Req any, Resp any](handler Handler[Req, Resp], opts ...Option) {
	c := defaultOpts()
	for _, o := range opts {
		o(&c)
	}

	lambda.Start(func(ctx context.Context, gatewayReq APIGatewayProxyRequest) (APIGatewayProxyResponse, error) {
		lambdaContext := Context[Req, Resp]{
			Context: ctx,
			Request: Request[Req]{
				req:        &gatewayReq,
				HTTPMethod: gatewayReq.HTTPMethod,
				Path:       gatewayReq.Path,
				PathParams: gatewayReq.PathParameters,
				Query:      Query(gatewayReq.MultiValueQueryStringParameters),
				Headers:    Headers(gatewayReq.MultiValueHeaders),
			},
			Response: Response[Resp]{
				StatusCode: http.StatusOK,
				Headers:    make(map[string][]string),
			},
			Locals: make(map[string]any),
		}

		// For the string -> []byte we need to use a more effective way. For now, let's keep the naive approach.
		err := populateLambdaContext(&gatewayReq, &lambdaContext)
		if err != nil {
			return c.errorHandler(err)
		}

		err = handler(&lambdaContext)
		if !(errors.Is(err, &lambdaContext.Response)) && err != nil {
			lambdaContext.error = err
		}
		if lambdaContext.Response.Err != nil {
			return c.errorHandler(lambdaContext.Response.Err)
		}

		r := APIGatewayProxyResponse{
			StatusCode:        lambdaContext.Response.StatusCode,
			MultiValueHeaders: lambdaContext.Response.Headers,
			Body:              lambdaContext.Response.Body.Bytes(),
		}
		return r, nil
	})
}

// populateLambdaContext will unmarshal the body from the gateway request into the lambda context.
func populateLambdaContext[Req any, Resp any](gatewayReq *APIGatewayProxyRequest, lambdaContext *Context[Req, Resp]) error {
	if gatewayReq.HTTPMethod == "GET" {
		return nil
	}
	if len(gatewayReq.Body) == 0 {
		return nil
	}
	var reader io.Reader = strings.NewReader(gatewayReq.Body)
	if gatewayReq.IsBase64Encoded {
		reader = base64.NewDecoder(base64.StdEncoding, reader)
	}
	return json.NewDecoder(reader).Decode(&lambdaContext.Request.Body)
}
