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
