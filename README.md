# errwrap

## Error Wrapper

Wrapping errors with some useful functions.

It uses generic so Go1.18+ needed.

## Description

`errwrap` provides some functions that wrapping, asserting, checking an error. These functions are generated by `Factory`, so the errors can be wrapped more systemically.

The reason of using generic is to distinguish errors. Errors must be wrapped in different methods, and wrapping by their types make classifying easier.

## BaseError

First of all, you can make a basic error with using `New` function. `errwrap` provides 4 types of `New` functions, which are `New`, `Newf`, `NewTyped`, `NewTypedf`.

```go
errwrap.New("new")
errwrap.Newf("%s%s", "new", "f")
errwrap.NewTyped[T]("newTyped")
errwrap.NewTypedf[T]("%s%s", "newTyped", "f")
```

The error generated by one of these functions is able to be formatted like wrappedError which will be explained from now on.

## Factory

Go `error` interface instances can be wrapped by this package's wrappedError. To wrap, you need to get a wrapper using a factory.

`Factory` needs a type as a generic which is representing where this error occured. And message is also needed, it explains what this error is.

```go
wrapper, assertor := Factory[UserController]("user controller")

// format specifier supplied
wrapper, assertor := Factoryf[UserController]("%s controller", "user")
```

`Factory` returns two functions which are `wrapper` and `assertor`. `wrapper` can wrap an error. The only thing you have to do to wrap is just inserting the error to the `wrapper`.

```go
err = wrapper(err)
```

When the error has a lot of wrapping depth, and when if you would like to pick a specified error that wrapped by this package, you can use `assertor` to do an assertion. Then `assertor` is going to return an asserted error and success or failure of the assertion. If it's failed, then the firstly returned argument will be `nil`.
To use `assertor` is same as `wrapper`.

```go
asserted, ok := assertor(err)

fmt.Println(asserted)  // user controller: ...
```

## Checker

`checker` works same as `errors.Is`, which checks a given error's type matching to the comparison's. To created this, `wrapper` is needed. `wrapper` becomes an argument of the function `WithChecker`, to create `checker`.

```go
checker := WithChecker(&wrapper)
```

After creating a `checker`, `wrapper` can make a checkable error. To use `checker` is like below:

```go
if checker(err) {
    fmt.Println("same")
}
```

## Format

`wrappedError` stores an error stack, so it can be printed in detail using `%+v` or `%+s`.

```bash
fmt.Printf("%+v\n", err)

# this will be shown like

github.com/p9595jh/errwrap.TestChecker
    /Users/medium/Desktop/PJH/blockchain/__xyz/2023.03/errwrap/errwrap_test.go:69
github.com/p9595jh/errwrap.TestChecker
    /Users/medium/Desktop/PJH/blockchain/__xyz/2023.03/errwrap/errwrap_test.go:68
github.com/p9595jh/errwrap.TestChecker
    /Users/medium/Desktop/PJH/blockchain/__xyz/2023.03/errwrap/errwrap_test.go:61
github.com/p9595jh/errwrap.TestChecker
    /Users/medium/Desktop/PJH/blockchain/__xyz/2023.03/errwrap/errwrap_test.go:60
```

## Benchmark

Comparison with `github.com/pkg/errors`.

Test code:

```go
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
```

Run:

```
$ go test -bench=. -benchtime=10s -benchmem -count 5

goos: darwin
goarch: amd64
pkg: github.com/p9595jh/errwrap
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkWrappedError-12          165193             72756 ns/op           28000 B/op        300 allocs/op
BenchmarkWrappedError-12          152595             72621 ns/op           28000 B/op        300 allocs/op
BenchmarkWrappedError-12          155252             73013 ns/op           28000 B/op        300 allocs/op
BenchmarkWrappedError-12          175780             72772 ns/op           28000 B/op        300 allocs/op
BenchmarkWrappedError-12          149059             71417 ns/op           28000 B/op        300 allocs/op
BenchmarkErrors-12                152004             82268 ns/op           33600 B/op        400 allocs/op
BenchmarkErrors-12                150874             78414 ns/op           33600 B/op        400 allocs/op
BenchmarkErrors-12                152422             79438 ns/op           33600 B/op        400 allocs/op
BenchmarkErrors-12                155005             85094 ns/op           33600 B/op        400 allocs/op
BenchmarkErrors-12                115872             88641 ns/op           33600 B/op        400 allocs/op
```

|                         | Wrappable in 10s |   ns/op |  B/op | allocs/op |
| ----------------------- | ---------------: | ------: | ----: | --------: |
| `errwrap`               |         159575.8 | 72515.8 | 28000 |       300 |
| `github.com/pkg/errors` |         145235.4 |   82771 | 33600 |       400 |

Processing ability of `errwrap` is improved around 9.87% than `github.com/pkg/errors`.
