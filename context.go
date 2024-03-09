package lambda

import (
	"context"
)

// None is an empty struct used when we are not interested on the request or response body.
type None struct{}

type Context[Req any, Resp any] struct {
	Context  context.Context
	Request  Request[Req]
	Response Response[Resp]
	Locals   map[string]any
	error    error
}

func (l *Context[Req, Resp]) Error() string {
	if l.error != nil {
		return l.error.Error()
	}
	return ""
}
