package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
)

type Handler[Req any, Resp any] func(*Context[Req, Resp]) error

// StartV1 will start the lambda function with the given handler and options.
//
// The handler is a function that receives a context and a pointer to a Context.
func StartV1[Req any, Resp any](handler Handler[Req, Resp], opts ...HttpOption) {
	c := defaultOpts()
	for _, o := range opts {
		o(&c)
	}

	startCtx := context.Background()
	for _, r := range c.resources {
		if err := r.Start(startCtx); err != nil {
			panic(fmt.Errorf("failed to start resource %s: %w", r.Name(), err))
		}
	}

	lambda.Start(func(ctx context.Context, gatewayReq APIGatewayProxyRequest) (APIGatewayProxyResponse, error) {
		req := Request[Req]{
			HTTPMethod: gatewayReq.HTTPMethod,
			Path:       gatewayReq.Path,
			PathParams: gatewayReq.PathParameters,
			Query:      Query(gatewayReq.QueryStringParameters),
			Headers:    Headers(gatewayReq.Headers),
			rawCookies: gatewayReq.MultiValueHeaders["Cookie"],
		}
		resp := Response[Resp]{
			StatusCode: http.StatusOK,
			Headers:    make(map[string]string),
		}

		lambdaContext := Context[Req, Resp]{
			Context:  ctx,
			Request:  &req,
			Response: &resp,
			Locals:   make(map[string]any),
		}

		// For the string -> []byte we need to use a more effective way. For now, let's keep the naive approach.
		err := populateLambdaContextV1(&gatewayReq, &lambdaContext)
		if err != nil {
			return toV1Response(c.errorHandler(ctx, err))
		}

		err = handler(&lambdaContext)
		if !(errors.Is(err, lambdaContext.Response)) && err != nil {
			lambdaContext.error = err
		}
		if lambdaContext.Response.Err != nil {
			return toV1Response(c.errorHandler(ctx, lambdaContext.Response.Err))
		}

		// lambdaContext.Response.Headers["Set-Cookie"] = strings.Join(lo, "; ")

		r := APIGatewayProxyResponse{
			StatusCode: lambdaContext.Response.StatusCode,
			Headers:    lambdaContext.Response.Headers,
			Body:       lambdaContext.Response.Body.Bytes(),
		}
		return r, nil
	})
}

func toV1Response(response HttpResponse, err error) (APIGatewayProxyResponse, error) {
	if err != nil {
		return APIGatewayProxyResponse{}, err
	}
	return APIGatewayProxyResponse(response), nil
}

// populateLambdaContextV1 will unmarshal the body from the gateway request into the lambda context.
func populateLambdaContextV1[Req any, Resp any](gatewayReq *APIGatewayProxyRequest, lambdaContext *Context[Req, Resp]) error {
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
