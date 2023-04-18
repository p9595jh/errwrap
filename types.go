package errwrap

type (
	NONE            struct{}
	Wrapper         func(error) error
	Checker         func(error) bool
	Assertor[T any] func(error) (*wrappedError[T], bool)
	unwrapper       interface{ Unwrap() error }
	counterer       interface{ counter() uintptr }
	updater         interface{ updatePC(int) }
)
