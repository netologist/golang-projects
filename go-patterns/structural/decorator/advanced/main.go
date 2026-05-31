package main

import (
	"context"
	"fmt"
	"os"
	"sync"
)

type inMemStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func newInMem() *inMemStore { return &inMemStore{data: map[string]string{}} }

func (s *inMemStore) Get(_ context.Context, key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

func (s *inMemStore) Set(_ context.Context, key, val string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = val
	return nil
}

func main() {
	ctx := context.Background()

	// Stack: Logging( Caching( InMem ) )
	store := NewLoggingStore(NewCachingStore(newInMem()), os.Stdout)

	store.Set(ctx, "lang", "go")
	store.Get(ctx, "lang") // cache miss — goes to inner
	store.Get(ctx, "lang") // cache hit — inner not called, logging still fires
	fmt.Println("done")
}
