package http

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithErrorHandler(t *testing.T) {
	o := options{}
	h := func(err error) (HttpResponse, error) {
		return HttpResponse{}, nil
	}
	WithErrorHandler(h)(&o)
	assert.Equal(t, fmt.Sprintf("%p", h), fmt.Sprintf("%p", o.errorHandler))
}

func TestWithResources(t *testing.T) {
	o := options{}
	r1 := &mockResource{"r1"}
	r2 := &mockResource{"r2"}
	WithResources(r1, r2)(&o)
	assert.Contains(t, o.resources, r1)
	assert.Contains(t, o.resources, r2)
}

type mockResource struct {
	value string
}

func (m *mockResource) Name() string {
	return m.value
}

func (m *mockResource) Start(ctx context.Context) error {
	return nil
}
