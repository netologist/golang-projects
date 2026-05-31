package main

import (
	"sync"
	"testing"
)

func TestGetDB_sameInstance(t *testing.T) {
	db1 := GetDB()
	db2 := GetDB()
	if db1 != db2 {
		t.Error("expected same instance, got different pointers")
	}
}

func TestGetDB_concurrentInit(t *testing.T) {
	Reset()
	var wg sync.WaitGroup
	instances := make([]*DB, 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			instances[i] = GetDB()
		}(i)
	}
	wg.Wait()
	for i, db := range instances {
		if db != instances[0] {
			t.Errorf("instance[%d] differs from instance[0]", i)
		}
	}
}
