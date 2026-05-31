package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func fetch(ctx context.Context, name string, fail bool) error {
	select {
	case <-time.After(10 * time.Millisecond):
		if fail {
			return fmt.Errorf("fetch %s failed", name)
		}
		fmt.Printf("fetched %s\n", name)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error { return fetch(ctx, "a", false) })
	g.Go(func() error { return fetch(ctx, "b", true) })
	g.Go(func() error { return fetch(ctx, "c", false) })

	if err := g.Wait(); err != nil {
		fmt.Println("group error:", err)
	}
}
