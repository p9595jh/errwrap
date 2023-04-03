package errwrap

import (
	"testing"

	"github.com/pkg/errors"
)

type (
	First  struct{}
	Second struct{}
	Third  struct{ I int }
	Fourth struct{ I int }
)

var (
	w1, a1 = Factory[First]("first")
	w2, a2 = Factory[Second]("second")
	w3, a3 = Factory[Third]("third")
	w4, a4 = Factory[Fourth]("fourth")

	_, _, _, _ = a1, a2, a3, a4
)

func TestWrapper(t *testing.T) {
	err := New("sample")
	err = w1(err)
	err = w2(err)
	err = w3(err)
	err = w4(err)
	_ = err

	t.Logf("%+v", err)
}

func TestErrors(t *testing.T) {
	err := errors.New("sample")
	err = errors.Wrap(err, "test1")
	err = errors.Wrap(err, "test2")
	t.Logf("%+v", err)
}

func TestAssertor(t *testing.T) {
	err := New("sample")
	err = w1(err)
	err = w2(err)
	err = w3(err)
	err = w4(err)

	t.Log(a2(err))
	t.Log(a3(err))

	_, an := Factory[struct{}]("none")
	t.Log(an(err))
}

func TestChecker(t *testing.T) {
	err := New("sample")
	err = w1(err)
	err = w2(err)

	c1 := WithChecker(&w1)
	c2 := WithChecker(&w2)
	c3 := WithChecker(&w3)
	c4 := WithChecker(&w4)

	err = w3(err)
	err = w4(err)

	t.Log(c1(err)) // true
	t.Log(c2(err)) // true
	t.Log(c3(err)) // false
	t.Log(c4(err)) // false

	t.Logf("%+v", err)
}

func BenchmarkWrappedError(b *testing.B) {
	w, _ := Factory[NONE]("test")
	err := New("sample")
	for i := 0; i < b.N; i++ {
		err := err
		for j := 0; j < 100; j++ {
			err = w(err)
		}
		_ = err
	}
}

func BenchmarkErrors(b *testing.B) {
	err := errors.New("sample")
	for i := 0; i < b.N; i++ {
		err := err
		for j := 0; j < 100; j++ {
			err = errors.Wrap(err, "test")
		}
		_ = err
	}
}
