package main

// Result holds either a value or an error.
type Result[T any] struct {
	val T
	err error
}

// Ok wraps a success value.
func Ok[T any](v T) Result[T] { return Result[T]{val: v} }

// Err wraps an error.
func Err[T any](e error) Result[T] { return Result[T]{err: e} }

// Unwrap returns the value and error.
func (r Result[T]) Unwrap() (T, error) { return r.val, r.err }

// Map applies fn if the Result is Ok.
func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.err != nil {
		return r
	}
	return Ok(fn(r.val))
}
