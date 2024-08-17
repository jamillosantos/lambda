package http

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type Query = mapUtils

type Headers = mapUtils

type PathParams map[string]string

type Request[T any] struct {
	HTTPMethod string
	Path       string
	PathParams PathParams
	Query      Query
	Headers    Headers
	rawCookies []string
	Body       T

	parseCookiesOnce sync.Once
	cookies          map[string]string
}

func (p PathParams) String(key string) (string, bool) {
	v, ok := p[key]
	return v, ok
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

func (r *Request[T]) Cookie(key string) (string, bool) {
	r.parseCookiesOnce.Do(func() {
		r.parseCookies()
	})
	v, ok := r.cookies[key]
	return v, ok
}

func (r *Request[T]) parseCookies() {
	// Cookie: name=value; name2=value2; name3=value3
	r.cookies = make(map[string]string, len(r.rawCookies))
	var (
		cookie string
		ok     bool
	)
	for _, cookiesH := range r.rawCookies {
		for {
			cookie, cookiesH, ok = strings.Cut(cookiesH, ";")
			if !ok && len(cookie) == 0 {
				break
			}
			key, cookie, ok := strings.Cut(cookie, "=")
			if !ok {
				continue
			}
			cookie, ok = extractCookieValue(cookie)
			if !ok {
				continue
			}
			r.cookies[strings.TrimLeft(key, " ")] = cookie
			if len(cookiesH) == 0 {
				break
			}
		}
	}
}

// TODO Validate cookie-octet value.
func extractCookieValue(cookie string) (string, bool) {
	l := len(cookie)
	if l > 0 && cookie[0] == '"' && l > 1 && cookie[len(cookie)-1] == '"' {
		return cookie[1 : len(cookie)-1], true
	}
	return cookie, true
}
