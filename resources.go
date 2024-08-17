package lambda

import "context"

// Resource is an interface that represents a resource that can be started and stopped.
// Example: A connection to the database or message broker, etc;
type Resource interface {
	Name() string
	Start(context.Context) error
}
