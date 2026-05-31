package main

// Result holds either a success value or an error.
type Result[T any] struct {
	val T
	err error
}

// Ok wraps a success value.
func Ok[T any](v T) Result[T] { return Result[T]{val: v} }

// Err wraps an error.
func Err[T any](e error) Result[T] { return Result[T]{err: e} }

// IsOk reports whether the Result holds a value.
func (r Result[T]) IsOk() bool { return r.err == nil }

// Unwrap returns the value and error.
func (r Result[T]) Unwrap() (T, error) { return r.val, r.err }

// Map applies fn if Ok, otherwise propagates the error.
func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.err != nil {
		return r
	}
	return Ok(fn(r.val))
}

// FlatMap applies fn (which itself returns a Result) if Ok.
func (r Result[T]) FlatMap(fn func(T) Result[T]) Result[T] {
	if r.err != nil {
		return r
	}
	return fn(r.val)
}

// Or returns the value if Ok, otherwise fallback.
func (r Result[T]) Or(fallback T) T {
	if r.err != nil {
		return fallback
	}
	return r.val
}
