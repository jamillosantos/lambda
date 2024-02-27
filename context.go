package lambda

import (
	"bytes"
	"encoding/json"
)

type Context[R any] struct {
	Request R
	status  int
	header  map[string][]string
	body    bytes.Buffer
	error   error
}

func (l *Context[R]) Status(status int) *Context[R] {
	l.status = status
	return l
}

func (l *Context[R]) Header(key string, value string) *Context[R] {
	if l.header == nil {
		l.header = make(map[string][]string)
	}
	h, ok := l.header[key]
	if !ok {
		l.header[key] = []string{value}
		return l
	}
	l.header[key] = append(h, value)
	return l
}

func (l *Context[R]) JSON(data any) *Context[R] {
	l.header["Content-Type"] = []string{"application/json"}
	if l.body.Len() > 0 {
		l.body.Reset()
	}
	err := json.NewEncoder(&l.body).Encode(data)
	if err != nil {
		l.error = err
	}
	return l
}

func (l *Context[R]) SendString(data string) *Context[R] {
	if l.body.Len() > 0 {
		l.body.Reset()
	}
	_, err := l.body.WriteString(data)
	if err != nil {
		l.error = err
	}
	return l
}

func (l *Context[R]) Error() string {
	if l.error != nil {
		return l.error.Error()
	}
	return ""
}
