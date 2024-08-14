package lambda

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type APIGatewayV2HTTPResponse struct {
	StatusCode        int                 `json:"statusCode"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	Body              string              `json:"body"`
	IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
	Cookies           []string            `json:"cookies"`
}

// StartV2 will start the lambda function with the given handler and options.
//
// The handler is a function that receives a context and a pointer to a Context.
func StartV2[Req any, Resp any](handler Handler[Req, Resp], opts ...Option) {
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

	lambda.Start(func(ctx context.Context, gatewayReq events.APIGatewayV2HTTPRequest) (APIGatewayV2HTTPResponse, error) {
		req := Request[Req]{
			HTTPMethod: gatewayReq.RequestContext.HTTP.Method,
			Path:       gatewayReq.RawPath,
			PathParams: gatewayReq.PathParameters,
			Query:      Query(gatewayReq.QueryStringParameters),
			Headers:    Headers(gatewayReq.Headers),
			rawCookies: gatewayReq.Cookies,
		}
		resp := Response[Resp]{
			StatusCode: http.StatusOK,
			Headers:    make(map[string]string),
			Cookies:    make([]Cookie, 0),
		}

		lambdaContext := Context[Req, Resp]{
			Context:  ctx,
			Request:  &req,
			Response: &resp,
			Locals:   make(map[string]any),
		}

		// For the string -> []byte we need to use a more effective way. For now, let's keep the naive approach.
		err := populateLambdaContextV2(&gatewayReq, &lambdaContext)
		if err != nil {
			return toV2Response(c.errorHandler(err))
		}

		err = handler(&lambdaContext)
		if !(errors.Is(err, lambdaContext.Response)) && err != nil {
			lambdaContext.error = err
		}
		if lambdaContext.Response.Err != nil {
			return toV2Response(c.errorHandler(lambdaContext.Response.Err))
		}

		r := APIGatewayV2HTTPResponse{
			StatusCode: lambdaContext.Response.StatusCode,
			Headers:    lambdaContext.Response.Headers,
			Cookies:    toCookieString(lambdaContext.Response.Cookies),
			Body:       lambdaContext.Response.Body.String(),
		}
		_ = json.NewEncoder(os.Stdout).Encode(&r)
		fmt.Println()
		return r, nil
	})
}

func toCookieString(cookies []Cookie) []string {
	cookieStrings := make([]string, len(cookies))
	for i, c := range cookies {
		cookieStrings[i] = c.String()
	}
	return cookieStrings
}

func toV2Response(response HTTPResponse, err error) (APIGatewayV2HTTPResponse, error) {
	if err != nil {
		return APIGatewayV2HTTPResponse{}, err
	}
	return APIGatewayV2HTTPResponse{
		StatusCode:      response.StatusCode,
		Headers:         response.Headers,
		Body:            string(response.Body),
		Cookies:         nil,
		IsBase64Encoded: false,
	}, nil
}

// populateLambdaContextV1 will unmarshal the body from the gateway request into the lambda context.
func populateLambdaContextV2[Req any, Resp any](gatewayReq *events.APIGatewayV2HTTPRequest, lambdaContext *Context[Req, Resp]) error {
	if gatewayReq.RequestContext.HTTP.Method == "GET" {
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
