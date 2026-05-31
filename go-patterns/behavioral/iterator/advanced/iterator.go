package main

import "iter"

// FromSlice adapts a slice to iter.Seq[T].
func FromSlice[T any](items []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range items {
			if !yield(v) {
				return
			}
		}
	}
}

// Filter yields only items satisfying pred.
func Filter[T any](seq iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if pred(v) && !yield(v) {
				return
			}
		}
	}
}

// Map transforms each item.
func Map[T, U any](seq iter.Seq[T], fn func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

// Take yields at most n items.
func Take[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		for v := range seq {
			if count >= n {
				return
			}
			if !yield(v) {
				return
			}
			count++
		}
	}
}
