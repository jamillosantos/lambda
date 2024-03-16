package lambda

import (
	"bytes"
	"encoding/json"
)

type Response[T any] struct {
	StatusCode int
	Headers    map[string][]string
	Body       bytes.Buffer
	Err        error
}

func (r *Response[T]) Redirect(url string, status int) {
	r.StatusCode = status
	r.Headers["Location"] = []string{url}
}

func (r *Response[T]) Status(status int) *Response[T] {
	r.StatusCode = status
	return r
}

func (r *Response[T]) Header(key string, value string) *Response[T] {
	h, ok := r.Headers[key]
	if !ok {
		r.Headers[key] = []string{value}
		return r
	}
	r.Headers[key] = append(h, value)
	return r
}

func (r *Response[T]) JSON(data any) *Response[T] {
	r.Headers["Content-Type"] = []string{"application/json"}
	if r.Body.Len() > 0 {
		r.Body.Reset()
	}
	err := json.NewEncoder(&r.Body).Encode(data)
	if err != nil {
		r.Err = err
	}
	return r
}

func (r *Response[T]) SendString(data string) *Response[T] {
	if r.Body.Len() > 0 {
		r.Body.Reset()
	}

	err := json.NewEncoder(&r.Body).Encode(data)
	if err != nil {
		r.Err = err
	}
	return r
}

func (r *Response[T]) Error() string {
	if r.Err != nil {
		return r.Err.Error()
	}
	return ""
}
