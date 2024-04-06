package lambda

import (
	"errors"
	"net/http"
	"testing"

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
			"key": []string{"value"},
		},
	}
	l := r.Header("key", "value2")
	assert.Equal(t, r, l)
	assert.Equal(t, "value2", r.Headers["key"][1])
}

func TestResponse_JSON(t *testing.T) {
	r := &Response[None]{
		Headers: Headers{},
	}
	l := r.JSON(map[string]any{"a": 1})
	assert.Equal(t, r, l)
	assert.Equal(t, "application/json", r.Headers["Content-Type"][0])
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
