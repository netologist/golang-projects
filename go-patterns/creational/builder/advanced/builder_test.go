package main

import (
	"errors"
	"testing"
)

func TestQueryBuilder_happyPath(t *testing.T) {
	q, err := NewQueryBuilder("orders").
		Select("id", "total").
		Where("status = 'pending'").
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := q.SQL()
	want := "SELECT id, total FROM orders WHERE status = 'pending'"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestQueryBuilder_missingTable(t *testing.T) {
	_, err := NewQueryBuilder("").Build()
	if !errors.Is(err, ErrMissingTable) {
		t.Errorf("expected ErrMissingTable, got %v", err)
	}
}

func TestQueryBuilder_multipleWhere(t *testing.T) {
	q, _ := NewQueryBuilder("items").
		Where("price > 10").
		Where("stock > 0").
		Build()
	got := q.SQL()
	want := "SELECT * FROM items WHERE price > 10 AND stock > 0"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
