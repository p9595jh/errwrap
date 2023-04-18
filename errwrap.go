package errwrap

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strconv"
)

// baseError represents pure error, so it does not prorivde wrapping.
// Other functions are same as WrappedError, e.g. StackTrace, Generic assertion.
type baseError[T any] struct {
	Message string
	pc      uintptr
}

func (e *baseError[T]) counter() uintptr {
	return e.pc
}

func (e *baseError[T]) Error() string {
	return e.Message
}

func (e *baseError[T]) Stack() []uintptr {
	return []uintptr{e.pc}
}

func (e *baseError[T]) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		switch {
		case s.Flag('+'):
			fn := runtime.FuncForPC(e.pc)
			file, line := fn.FileLine(e.pc)
			name := fn.Name()
			io.WriteString(s, e.Error())
			io.WriteString(s, "\n")
			io.WriteString(s, name)
			io.WriteString(s, "\n\t")
			io.WriteString(s, file)
			io.WriteString(s, ":")
			io.WriteString(s, strconv.Itoa(line))
		default:
			io.WriteString(s, e.Error())
		}
	default:
		io.WriteString(s, "{%!")
		io.WriteString(s, string(verb))
		io.WriteString(s, "(error=")
		io.WriteString(s, e.Error())
		io.WriteString(s, ")}")
	}
}

// wrappedError wraps an error and stores its stacks.
// It also has a generic to be used in assertion.
// Generic is usually set as a component where it's called.
type wrappedError[T any] struct {
	Inner   error
	Message string
	pc      uintptr
}

func (e *wrappedError[T]) counter() uintptr {
	return e.pc
}

func (e *wrappedError[T]) Error() string {
	if e.Inner == nil {
		return e.Message
	}
	return e.Message + ": " + e.Inner.Error()
}

func (e *wrappedError[T]) Unwrap() error { return e.Inner }

// Overrides Format to format itself.
//
// %+v and %+s print stack in detail.
func (e *wrappedError[T]) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		switch {
		case s.Flag('+'):
			io.WriteString(s, e.Error())
			err := error(e)
			for {
				// Format detail if it's an instance of wrappedError.
				// Printed shape is like below:
				//
				// [package].[function] ([error])
				//     /[path]/[file].go:[line]
				//
				// If it's not an instance of wrappedError, then it's like:
				//
				// [untraceable]
				//     [error]
				if counterable, ok := err.(counterer); ok {
					pc := counterable.counter()
					fn := runtime.FuncForPC(pc)
					file, line := fn.FileLine(pc)
					name := fn.Name()
					io.WriteString(s, "\n")
					io.WriteString(s, name)
					io.WriteString(s, " (")
					io.WriteString(s, err.Error())
					io.WriteString(s, ")\n\t")
					io.WriteString(s, file)
					io.WriteString(s, ":")
					io.WriteString(s, strconv.Itoa(line))
				} else {
					io.WriteString(s, "\n")
					io.WriteString(s, untraceable)
					io.WriteString(s, "\n\t")
					io.WriteString(s, err.Error())
				}

				// It can keep formatting after untraceable error if it's unwrappable.
				// On the other hand, it only can work when the latest error is an instance
				// of wrappedError.
				if unwrappable, ok := err.(unwrapper); ok {
					err = unwrappable.Unwrap()
				} else {
					break
				}
			}
		default:
			io.WriteString(s, e.Error())
		}
	default:
		io.WriteString(s, "{%!")
		io.WriteString(s, string(verb))
		io.WriteString(s, "(error=")
		io.WriteString(s, e.Error())
		io.WriteString(s, ")}")
	}
}

// updatePC updates the last stack to caller with the given skip.
func (e *wrappedError[T]) updatePC(skip int) {
	e.pc = caller(skip)
}

// factory returns two functions which are a wrapper and an assertor.
// wrapper can wrap an error and append the error stack.
// assertor does assertion if it is possible.
func factory[T any](message string) (w Wrapper, a Assertor[T]) {
	w = func(err error) error {
		if err == nil {
			return nil
		}
		return &wrappedError[T]{
			Message: message,
			Inner:   err,
			pc:      caller(2),
		}
	}

	a = func(err error) (*wrappedError[T], bool) {
		var (
			target *wrappedError[T]
			ok     = errors.As(err, &target)
		)
		return target, ok
	}

	return
}

// Factory must be typed with a generic and gets a message.
// This returns a wrapper and an assertor that can be used for wrapping and an assertion.
func Factory[T any](message string) (Wrapper, Assertor[T]) {
	return factory[T](message)
}

// Factoryf is same as Factory and the only difference is that it has a format specifier for the message.
func Factoryf[T any](message string, args ...any) (Wrapper, Assertor[T]) {
	return factory[T](fmt.Sprintf(message, args...))
}

// WithChecker attaches a checker to the given wrapper.
// Checker is able to check the given error is the type of the wrapper.
func WithChecker(wrapper *Wrapper, size ...int) Checker {
	var sz int
	if len(size) > 0 {
		sz = size[0]
	}

	var (
		errs     = make(map[error]bool, sz)
		_wrapper = *wrapper // store previous
	)

	*wrapper = func(err error) error {
		if err == nil {
			return nil
		}
		werr := _wrapper(err)
		werr.(updater).updatePC(3)
		errs[Innermost(err)] = true
		return werr
	}

	return func(err error) bool {
		for k := range errs {
			if errors.Is(err, k) {
				return true
			}
		}
		return false
	}
}
