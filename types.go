package errwrap

type (
	NONE            struct{}
	wrapper         func(error) error
	checker         func(error) bool
	assertor[T any] func(error) (*wrappedError[T], bool)
	unwrapper       interface{ Unwrap() error }
	counterer       interface{ counter() uintptr }
	updater         interface{ updatePC(int) }
)
