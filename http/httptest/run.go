package httptest

import (
	"context"
	"encoding/json"
	"errors"
	lambdahttp "github.com/jamillosantos/lambda/http"
	"net/http"
)

type TestHttpContext[Req any, Resp any] struct {
	lambdahttp.Context[Req, Resp]
}

func (t *TestHttpContext[Req, Resp]) ResponseBody() (Resp, error) {
	var r Resp
	err := json.NewDecoder(&t.Response.Body).Decode(&r)
	if err != nil {
		return r, err
	}
	return r, nil
}

func Run[Req any, Resp any](handler lambdahttp.Handler[Req, Resp], opts ...Option) (*TestHttpContext[Req, Resp], error) {
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
	ctx := &TestHttpContext[Req, Resp]{
		lambdahttp.Context[Req, Resp]{
			Context: context.Background(),
			Request: &lambdahttp.Request[Req]{
				HTTPMethod: o.httpMethod,
				Path:       o.path,
				PathParams: o.pathParams,
				Query:      lambdahttp.Query(o.query),
				Headers:    lambdahttp.Headers(o.headers),
				Body:       o.req.(Req),
			},
			Response: &lambdahttp.Response[Resp]{
				StatusCode: http.StatusOK,
				Headers:    make(map[string]string),
			},
			Locals: o.locals,
		},
	}
	err := handler(&ctx.Context)
	if errors.Is(err, ctx.Response) { // nolint
		if ctx.Response.Err != nil {
			return nil, ctx.Response.Err
		}
	} else if err != nil {
		return nil, err
	}
	return ctx, nil
}
