package http

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResponse_Error(t *testing.T) {
	r := &Response[None]{
		Err: errors.New("error"),
	}
	assert.Equal(t, r.Error(), r.Err.Error())
}

func TestResponse_Header(t *testing.T) {
	r := &Response[None]{
		Headers: Headers{
			"key": "value",
		},
	}
	l := r.Header("key", "value2")
	assert.Equal(t, r, l)
	assert.Equal(t, "value2", r.Headers["key"])
}

func TestResponse_JSON(t *testing.T) {
	r := &Response[None]{
		Headers: Headers{},
	}
	l := r.JSON(map[string]any{"a": 1})
	assert.Equal(t, r, l)
	assert.Equal(t, "application/json", r.Headers["Content-Type"])
	assert.Equal(t, "{\"a\":1}\n", r.Body.String())
}

func TestResponse_Redirect(t *testing.T) {
	t.Run("should set the status code to 307 when no status is provided", func(t *testing.T) {
		r := &Response[None]{
			Headers: Headers{},
		}
		l := r.Redirect("http://example.com")
		assert.Equal(t, r, l)
		assert.Equal(t, http.StatusTemporaryRedirect, r.StatusCode)
	})

	t.Run("should set the status code to the provided status", func(t *testing.T) {
		r := &Response[None]{
			Headers: Headers{},
		}
		l := r.Redirect("http://example.com", http.StatusMovedPermanently)
		assert.Equal(t, r, l)
		assert.Equal(t, http.StatusMovedPermanently, r.StatusCode)
	})
}

func TestResponse_SendString(t *testing.T) {
	r := &Response[None]{}
	l := r.SendString("data")
	assert.Equal(t, r, l)
	assert.Equal(t, "\"data\"\n", r.Body.String())
}

func TestResponse_Status(t *testing.T) {
	r := &Response[None]{}
	l := r.Status(http.StatusAccepted)
	assert.Equal(t, r, l)
	assert.Equal(t, http.StatusAccepted, r.StatusCode)
}

func TestResponse_SetCookie(t *testing.T) {
	r := &Response[None]{
		Headers: Headers{},
	}
	wantCookie := Cookie{
		Name:  "name",
		Value: "value",
	}
	_ = r.SetCookie(wantCookie)
	assert.Equal(t, r.Cookies[0], wantCookie)
}

func TestCookie_String(t *testing.T) {
	t.Run("should set the cookie", func(t *testing.T) {
		givenCookie := Cookie{
			Name:  "name",
			Value: "value",
		}
		assert.Equal(t, "name=value", givenCookie.String())
	})

	t.Run("should set the cookie with path", func(t *testing.T) {
		givenCookie := Cookie{
			Name:  "name",
			Value: "value",
			Path:  "/",
		}
		assert.Equal(t, "name=value; path=/", givenCookie.String())
	})

	t.Run("should set the cookie with domain", func(t *testing.T) {
		givenCookie := Cookie{
			Name:   "name",
			Value:  "value",
			Domain: "example.com",
		}
		assert.Equal(t, "name=value; domain=example.com", givenCookie.String())
	})

	t.Run("should set the cookie with max age", func(t *testing.T) {
		givenCookie := Cookie{
			Name:   "name",
			Value:  "value",
			MaxAge: 1,
		}
		assert.Equal(t, "name=value; max-age=1", givenCookie.String())
	})

	t.Run("should set the cookie with expires", func(t *testing.T) {
		expires := time.Date(2024, 1, 2, 3, 4, 5, 6, time.UTC)
		givenCookie := Cookie{
			Name:    "name",
			Value:   "value",
			Expires: expires,
		}
		assert.Equal(t, "name=value; expires=Tue, 02 Jan 2024 03:04:05 UTC", givenCookie.String())
	})

	t.Run("should set the cookie with secure", func(t *testing.T) {
		givenCookie := Cookie{
			Name:   "name",
			Value:  "value",
			Secure: true,
		}
		assert.Equal(t, "name=value; secure", givenCookie.String())
	})

	t.Run("should set the cookie with http only", func(t *testing.T) {
		givenCookie := Cookie{
			Name:     "name",
			Value:    "value",
			HTTPOnly: true,
		}
		assert.Equal(t, "name=value; HttpOnly", givenCookie.String())
	})

	t.Run("should set the cookie with same site", func(t *testing.T) {
		// table tests with all possible same sites
		tests := []struct {
			name     string
			sameSite CookieSameSite
			expected string
		}{
			{"default", CookieSameSiteDefaultMode, "name=value; SameSite"},
			{"lax", CookieSameSiteLaxMode, "name=value; SameSite=Lax"},
			{"strict", CookieSameSiteStrictMode, "name=value; SameSite=Strict"},
			{"none", CookieSameSiteNoneMode, "name=value; SameSite=None"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				givenCookie := Cookie{
					Name:     "name",
					Value:    "value",
					SameSite: tt.sameSite,
				}
				assert.Equal(t, tt.expected, givenCookie.String())
			})
		}
	})
}
