package lambda

import "context"

// None is an empty struct used when we are not interested on the request or response body.
type None struct{}

type Context[Req any] struct {
	Context context.Context
	Request Req
	Locals  map[string]any
}

func (l *Context[Req]) SetLocal(key string, value any) *Context[Req] {
	l.Locals[key] = value
	return l
}

func (l *Context[Req]) GetLocal(key string) (any, bool) {
	value, ok := l.Locals[key]
	return value, ok
}

func (l *Context[Req]) UnsetLocal(key string) *Context[Req] {
	delete(l.Locals, key)
	return l
}
