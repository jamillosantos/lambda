package http

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContext_Error(t *testing.T) {
	t.Run("should return the error message", func(t *testing.T) {
		l := Context[None, None]{error: errors.New("error message")}
		require.Equal(t, "error message", l.Error())
	})

	t.Run("should return an empty string if there is no error", func(t *testing.T) {
		l := &Context[None, None]{}
		require.Equal(t, "", l.Error())
	})
}

func TestContext_GetLocal(t *testing.T) {
	l := &Context[None, None]{Locals: map[string]any{"key": "value"}}
	value, ok := l.GetLocal("key")
	require.True(t, ok)
	require.Equal(t, "value", value)
}

func TestContext_SetLocal(t *testing.T) {
	l := &Context[None, None]{Locals: map[string]any{}}
	assert.Equal(t, l.SetLocal("key", "value"), l)
	assert.Equal(t, map[string]any{"key": "value"}, l.Locals)
}

func TestContext_UnsetLocal(t *testing.T) {
	l := &Context[None, None]{Locals: map[string]any{"key": "value"}}
	assert.Equal(t, l.UnsetLocal("key"), l)
	assert.Equal(t, map[string]any{}, l.Locals)
}
