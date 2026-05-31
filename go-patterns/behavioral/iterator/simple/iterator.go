package main

// SliceIter returns a channel that yields items from the slice, then closes.
func SliceIter[T any](items []T) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, v := range items {
			ch <- v
		}
	}()
	return ch
}
