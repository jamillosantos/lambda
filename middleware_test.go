package lambda

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func mBuilder() (*[]string, func(step string) Middleware[None, None]) {
	seq := []string{}
	m := func(step string) Middleware[None, None] {
		return func(ctx *Context[None, None], next Handler[None, None]) error {
			seq = append(seq, step+"-before")
			err := next(ctx)
			seq = append(seq, step+"-after")
			return err
		}
	}
	return &seq, m
}

func TestUse(t *testing.T) {
	t.Run("should call the middleware and the handler", func(t *testing.T) {
		seq, m := mBuilder()
		err := Use(func(ctx *Context[None, None]) error {
			*seq = append(*seq, "handler")
			return nil
		}, m("outer"), m("inner"))(nil)
		require.NoError(t, err)
		require.Equal(t, []string{"outer-before", "inner-before", "handler", "inner-after", "outer-after"}, *seq)
	})

	t.Run("should interrupt middleware if next is not called", func(t *testing.T) {
		seq, m := mBuilder()
		err := Use(func(ctx *Context[None, None]) error {
			*seq = append(*seq, "handler")
			return nil
		}, m("outer"), func(ctx *Context[None, None], next Handler[None, None]) error {
			*seq = append(*seq, "interrupt")
			return nil
		})(nil)
		require.NoError(t, err)
		require.Equal(t, []string{"outer-before", "interrupt", "outer-after"}, *seq)
	})
}
