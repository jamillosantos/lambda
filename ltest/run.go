package ltest

import (
	"context"
	"net/http"

	"github.com/jamillosantos/lambda"
)

func Run[Req any, Resp any](handler lambda.Handler[Req, Resp], opts ...Option) (*lambda.Context[Req, Resp], error) {
	o := options{
		httpMethod: "GET",
		path:       "/",
		pathParams: make(map[string]string),
		query:      make(map[string][]string),
		headers:    make(map[string][]string),
	}
	var req Req
	o.req = req
	for _, opt := range opts {
		opt(&o)
	}
	ctx := &lambda.Context[Req, Resp]{
		Context: context.Background(),
		Request: lambda.Request[Req]{
			HTTPMethod: o.httpMethod,
			Path:       o.path,
			PathParams: o.pathParams,
			Query:      lambda.NewQuery(make(map[string]string), o.query),
			Headers:    lambda.NewHeaders(make(map[string]string), o.headers),
			Body:       o.req.(Req),
		},
		Response: lambda.Response[Resp]{
			StatusCode: http.StatusOK,
			Headers:    make(map[string][]string),
		},
		Locals: o.locals,
	}
	err := handler(ctx)
	if err == &ctx.Response { // nolint
		if ctx.Response.Err != nil {
			return nil, ctx.Response.Err
		}
	} else if err != nil {
		return nil, err
	}
	return ctx, nil
}
