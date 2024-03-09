package lambda

import (
	"bytes"
)

type Response[T any] struct {
	status int
	header map[string][]string
	body   bytes.Buffer
	err    error
	Body   T
}
