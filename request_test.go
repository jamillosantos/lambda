package lambda

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathParams_String(t *testing.T) {
	p := PathParams{
		"key": "value",
	}
	v, ok := p.String("key")
	assert.True(t, ok)
	assert.Equal(t, "value", v)
}

func TestPathParams_Bool(t *testing.T) {
	p := PathParams{
		"key":  "true",
		"key2": "invalid",
	}
	t.Run("should return true when the key exists and is valid", func(t *testing.T) {
		v, err := p.Bool("key")
		assert.NoError(t, err)
		assert.True(t, v)
	})

	t.Run("should fail when the key does not exist", func(t *testing.T) {
		_, err := p.Bool("key3")
		assert.ErrorIs(t, err, ErrKeyNotFound)
	})

	t.Run("should fail when the key exists but is invalid", func(t *testing.T) {
		_, err := p.Bool("key2")
		assert.Error(t, err)
	})
}

func TestPathParams_Float64(t *testing.T) {
	p := PathParams{
		"key":  "1.0",
		"key2": "invalid",
	}

	t.Run("should return the float64 value when the key exists and is valid", func(t *testing.T) {
		v, err := p.Float64("key")
		assert.NoError(t, err)
		assert.Equal(t, 1.0, v)
	})

	t.Run("should fail when the key does not exist", func(t *testing.T) {
		_, err := p.Float64("key3")
		assert.ErrorIs(t, err, ErrKeyNotFound)
	})

	t.Run("should fail when the key exists but is invalid", func(t *testing.T) {
		_, err := p.Float64("key2")
		assert.Error(t, err)
	})
}

func TestPathParams_Int(t *testing.T) {
	p := PathParams{
		"key":  "1",
		"key2": "invalid",
	}

	t.Run("should return the int value when the key exists and is valid", func(t *testing.T) {
		v, err := p.Int("key")
		assert.NoError(t, err)
		assert.Equal(t, 1, v)
	})

	t.Run("should fail when the key does not exist", func(t *testing.T) {
		_, err := p.Int("key3")
		assert.ErrorIs(t, err, ErrKeyNotFound)
	})

	t.Run("should fail when the key exists but is invalid", func(t *testing.T) {
		_, err := p.Int("key2")
		assert.Error(t, err)
	})
}

func TestPathParams_Int64(t *testing.T) {
	p := PathParams{
		"key":  "1",
		"key2": "invalid",
	}

	t.Run("should return the int64 value when the key exists and is valid", func(t *testing.T) {
		v, err := p.Int64("key")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), v)
	})

	t.Run("should fail when the key does not exist", func(t *testing.T) {
		_, err := p.Int64("key3")
		assert.ErrorIs(t, err, ErrKeyNotFound)
	})

	t.Run("should fail when the key exists but is invalid", func(t *testing.T) {
		_, err := p.Int64("key2")
		assert.Error(t, err)
	})
}

func TestRequest_Cookie(t *testing.T) {
	type testCase struct {
		name      string
		cookies   []string
		key       string
		wantValue string
		wantOk    bool
	}
	tests := []testCase{
		{
			name:      "should return the cookie value when it exists",
			cookies:   []string{"key=value"},
			key:       "key",
			wantValue: "value",
			wantOk:    true,
		},
		{
			name:    "should return the cookie value when it exists",
			cookies: []string{"key=value"},
			key:     "non-existing-key",
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request[None]{
				rawCookies: tt.cookies,
			}
			gotCookieValue, ok := r.Cookie(tt.key)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantValue, gotCookieValue)
		})
	}
}

func TestRequest_parseCookies(t *testing.T) {
	r := &Request[None]{
		rawCookies: []string{
			`key=value; key2="value2"; key3; key4=value 4`,
			`key5; key6="value6"`,
		},
	}
	r.parseCookies()

	assert.Len(t, r.cookies, 4)

	assert.Contains(t, r.cookies, "key")
	assert.Contains(t, r.cookies["key"], "value")

	assert.Contains(t, r.cookies, "key2")
	assert.Contains(t, r.cookies["key2"], "value2")

	assert.NotContains(t, r.cookies, "key3")

	assert.Contains(t, r.cookies, "key4")
	assert.Contains(t, r.cookies["key4"], "value 4")

	assert.NotContains(t, r.cookies, "key5")

	assert.Contains(t, r.cookies, "key6")
	assert.Contains(t, r.cookies["key6"], "value6")
}
