package errwrap

import "fmt"

func newBaseError[T any](message string) *baseError[T] {
	return &baseError[T]{
		Message: message,
		pc:      caller(3),
	}
}

// New creates a new error.
func New(message string) error {
	return newBaseError[NONE](message)
}

// Newf creates a new error with a format specifier.
func Newf(message string, args ...any) error {
	return newBaseError[NONE](fmt.Sprintf(message, args...))
}

// NewTyped creates a new typed error.
func NewTyped[T any](message string) error {
	return newBaseError[T](message)
}

// NewTypedf creates a new typed error with a format specifier.
func NewTypedf[T any](message string, args ...any) error {
	return newBaseError[T](fmt.Sprintf(message, args...))
}
