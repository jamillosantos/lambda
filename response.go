package lambda

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Response[T any] struct {
	StatusCode int
	Headers    map[string]string
	Cookies    []Cookie
	Body       bytes.Buffer
	Err        error
}

func (r *Response[T]) Redirect(url string, status ...int) *Response[T] {
	st := http.StatusTemporaryRedirect
	if len(status) > 0 {
		st = status[0]
	}
	r.StatusCode = st
	r.Headers["Location"] = url
	return r
}

func (r *Response[T]) Status(status int) *Response[T] {
	r.StatusCode = status
	return r
}

func (r *Response[T]) Header(key string, value string) *Response[T] {
	r.Headers[key] = value
	return r
}

func (r *Response[T]) JSON(data any) *Response[T] {
	r.Headers["Content-Type"] = "application/json"
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

type CookieSameSite string

const (
	CookieSameSiteDefaultMode CookieSameSite = "SameSite"
	CookieSameSiteLaxMode     CookieSameSite = "Lax"
	CookieSameSiteStrictMode  CookieSameSite = "Strict"
	CookieSameSiteNoneMode    CookieSameSite = "None"
)

type Cookie struct {
	Name        string
	Value       string
	Path        string
	Domain      string
	MaxAge      int
	Expires     time.Time
	Secure      bool
	HTTPOnly    bool
	SameSite    CookieSameSite
	SessionOnly bool
}

func (r *Response[T]) SetCookie(c Cookie) *Response[T] {
	if r.Cookies == nil {
		r.Cookies = make([]Cookie, 0)
	}
	r.Cookies = append(r.Cookies, c)
	return r
}

var (
	sbPool = sync.Pool{
		New: func() any {
			return &strings.Builder{}
		},
	}
)

func (c *Cookie) String() string {
	sb := sbPool.Get().(*strings.Builder)
	sb.Reset()
	defer sbPool.Put(sb)
	sb.WriteString(c.Name)
	sb.WriteString("=")
	sb.WriteString(c.Value)

	if c.MaxAge > 0 {
		sb.WriteString("; max-age=")
		sb.WriteString(strconv.Itoa(c.MaxAge))
	} else if !c.Expires.IsZero() {
		sb.WriteString("; expires=")
		sb.WriteString(c.Expires.UTC().Format(time.RFC1123))
	}
	if len(c.Domain) > 0 {
		sb.WriteString("; domain=")
		sb.WriteString(c.Domain)
	}
	if len(c.Path) > 0 {
		sb.WriteString("; path=")
		sb.WriteString(c.Path)
	}
	if c.HTTPOnly {
		sb.WriteString("; HttpOnly")
	}
	if c.Secure {
		sb.WriteString("; secure")
	}
	switch c.SameSite {
	case CookieSameSiteDefaultMode:
		sb.WriteString("; SameSite")
	case CookieSameSiteLaxMode, CookieSameSiteStrictMode, CookieSameSiteNoneMode:
		sb.WriteString("; SameSite=")
		sb.WriteString(string(c.SameSite))
	}
	return sb.String()
}

var unsetCookieDate = time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

func (r *Response[T]) UnsetCookie(name string) *Response[T] {
	return r.SetCookie(Cookie{
		Name:    name,
		Expires: unsetCookieDate,
	})
}
