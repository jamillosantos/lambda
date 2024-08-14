package ltest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jamillosantos/lambda"
)

type TestContext[Req any, Resp any] struct {
	lambda.Context[Req, Resp]
}

func (t *TestContext[Req, Resp]) ResponseBody() (Resp, error) {
	var r Resp
	err := json.NewDecoder(&t.Response.Body).Decode(&r)
	if err != nil {
		return r, err
	}
	return r, nil
}

func Run[Req any, Resp any](handler lambda.Handler[Req, Resp], opts ...Option) (*TestContext[Req, Resp], error) {
	o := options{
		httpMethod: "GET",
		path:       "/",
		pathParams: make(map[string]string),
		query:      make(map[string]string),
		headers:    make(map[string]string),
	}
	var req Req
	o.req = req
	for _, opt := range opts {
		opt(&o)
	}
	ctx := &TestContext[Req, Resp]{
		lambda.Context[Req, Resp]{
			Context: context.Background(),
			Request: &lambda.Request[Req]{
				HTTPMethod: o.httpMethod,
				Path:       o.path,
				PathParams: o.pathParams,
				Query:      lambda.Query(o.query),
				Headers:    lambda.Headers(o.headers),
				Body:       o.req.(Req),
			},
			Response: &lambda.Response[Resp]{
				StatusCode: http.StatusOK,
				Headers:    make(map[string]string),
			},
			Locals: o.locals,
		},
	}
	err := handler(&ctx.Context)
	if err == ctx.Response { // nolint
		if ctx.Response.Err != nil {
			return nil, ctx.Response.Err
		}
	} else if err != nil {
		return nil, err
	}
	return ctx, nil
}
