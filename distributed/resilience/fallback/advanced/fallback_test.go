package main

import (
	"context"
	"errors"
	"testing"
)

func TestFallback_PrimarySuccess(t *testing.T) {
	f := Fallback[string]{
		Primary:  func(_ context.Context) (string, error) { return "primary", nil },
		Fallback: func(_ context.Context) (string, error) { return "fallback", nil },
	}
	v, err := f.Execute(context.Background())
	if err != nil || v != "primary" {
		t.Fatalf("want primary, got %q %v", v, err)
	}
}

func TestFallback_PrimaryFail(t *testing.T) {
	f := Fallback[string]{
		Primary:  func(_ context.Context) (string, error) { return "", errors.New("down") },
		Fallback: func(_ context.Context) (string, error) { return "stale", nil },
	}
	v, err := f.Execute(context.Background())
	if err != nil || v != "stale" {
		t.Fatalf("want stale, got %q %v", v, err)
	}
}

func TestChain_AllFail(t *testing.T) {
	fail := func(_ context.Context) (string, error) { return "", errors.New("fail") }
	_, err := Chain(context.Background(), fail, fail, fail)
	if !errors.Is(err, ErrAllFailed) {
		t.Fatalf("want ErrAllFailed, got %v", err)
	}
}

func TestChain_MiddleSucceeds(t *testing.T) {
	fail := func(_ context.Context) (string, error) { return "", errors.New("fail") }
	ok := func(_ context.Context) (string, error) { return "ok", nil }
	v, err := Chain(context.Background(), fail, ok, fail)
	if err != nil || v != "ok" {
		t.Fatalf("want ok, got %q %v", v, err)
	}
}
