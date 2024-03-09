package lambda

import (
	"bytes"
	"encoding/json"
)

type Response[T any] struct {
	status int
	header map[string][]string
	body   bytes.Buffer
	err    error
	Body   T
}

func (r *Response[T]) Redirect(url string, status int) {
	r.status = status
	r.header["Location"] = []string{url}
}

func (r *Response[T]) Status(status int) *Response[T] {
	r.status = status
	return r
}

func (r *Response[T]) Header(key string, value string) *Response[T] {
	if r.header == nil {
		r.header = make(map[string][]string)
	}
	h, ok := r.header[key]
	if !ok {
		r.header[key] = []string{value}
		return r
	}
	r.header[key] = append(h, value)
	return r
}

func (r *Response[T]) JSON(data any) *Response[T] {
	r.header["Content-Type"] = []string{"application/json"}
	if r.body.Len() > 0 {
		r.body.Reset()
	}
	err := json.NewEncoder(&r.body).Encode(data)
	if err != nil {
		r.err = err
	}
	return r
}

func (r *Response[T]) SendString(data string) *Response[T] {
	if r.body.Len() > 0 {
		r.body.Reset()
	}

	err := json.NewEncoder(&r.body).Encode(data)
	if err != nil {
		r.err = err
	}
	return r
}

func (r *Response[T]) Error() string {
	if r.err != nil {
		return r.err.Error()
	}
	return ""
}
