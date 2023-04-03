package errwrap

import "runtime"

const untraceable = "UNTRACEABLE"

// caller calls the function's pointer to figure out who has called.
// It skips the stack as specified.
func caller(skip int) uintptr {
	pc, _, _, _ := runtime.Caller(skip)
	return pc
}

// Innermost returns the innermost error of the given error.
func Innermost(err error) error {
	for werr, ok := err.(unwrapper); ok; werr, ok = err.(unwrapper) {
		err = werr.Unwrap()
	}
	return err
}
