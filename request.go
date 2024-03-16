package lambda

import (
	"fmt"
	"strconv"
)

type Query = mapArrayUtils

type Headers = mapArrayUtils

type PathParams map[string]string

type Request[T any] struct {
	req        *APIGatewayProxyRequest
	HTTPMethod string
	Path       string
	PathParams PathParams
	Query      Query
	Headers    Headers
	Body       T
}

func (p PathParams) String(key string) string {
	return p[key]
}

func (p PathParams) Int(key string) (int, error) {
	v, ok := p[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	return strconv.Atoi(v)
}

func (p PathParams) Int64(key string) (int64, error) {
	v, ok := p[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	return strconv.ParseInt(v, 10, 64)
}

func (p PathParams) Float64(key string) (float64, error) {
	v, ok := p[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	return strconv.ParseFloat(v, 64)
}

func (p PathParams) Bool(key string) (bool, error) {
	v, ok := p[key]
	if !ok {
		return false, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	return strconv.ParseBool(v)
}
