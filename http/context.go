package http

import (
	"context"
	"errors"
)

type Context[Req any, Resp any] struct {
	Context  context.Context
	Request  *Request[Req]
	Response *Response[Resp]
	Locals   map[string]any
	error    error
}

func (l *Context[Req, Resp]) Error() string {
	if l.error != nil {
		return l.error.Error()
	}
	return ""
}

func (l *Context[Req, Resp]) SetLocal(key string, value any) *Context[Req, Resp] {
	l.Locals[key] = value
	return l
}

func (l *Context[Req, Resp]) GetLocal(key string) (any, bool) {
	value, ok := l.Locals[key]
	return value, ok
}

func (l *Context[Req, Resp]) UnsetLocal(key string) *Context[Req, Resp] {
	delete(l.Locals, key)
	return l
}

func (l *Context[Req, Resp]) Is(err error) bool {
	return errors.Is(l.error, err)
}
