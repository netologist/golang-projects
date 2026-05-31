package main

import (
	"errors"
	"fmt"
	"strings"
)

// ErrMissingTable is returned when Build is called without a table.
var ErrMissingTable = errors.New("table is required")

// Query is an immutable, validated SQL query.
type Query struct {
	sql  string
	args []any
}

// SQL returns the rendered SQL string.
func (q Query) SQL() string { return q.sql }

// Args returns the bound arguments.
func (q Query) Args() []any { return q.args }

// QueryBuilder accumulates clauses and produces a Query via Build().
type QueryBuilder struct {
	table  string
	cols   []string
	wheres []string
	args   []any
}

// NewQueryBuilder starts a builder for the given table.
func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{table: table}
}

// Select appends columns to select.
func (b *QueryBuilder) Select(cols ...string) *QueryBuilder {
	b.cols = append(b.cols, cols...)
	return b
}

// Where appends a condition (joined with AND) and its bound args.
func (b *QueryBuilder) Where(cond string, args ...any) *QueryBuilder {
	b.wheres = append(b.wheres, cond)
	b.args = append(b.args, args...)
	return b
}

// Build validates and renders the immutable Query.
func (b *QueryBuilder) Build() (Query, error) {
	if b.table == "" {
		return Query{}, fmt.Errorf("build query: %w", ErrMissingTable)
	}
	cols := "*"
	if len(b.cols) > 0 {
		cols = strings.Join(b.cols, ", ")
	}
	sql := fmt.Sprintf("SELECT %s FROM %s", cols, b.table)
	if len(b.wheres) > 0 {
		sql += " WHERE " + strings.Join(b.wheres, " AND ")
	}
	return Query{sql: sql, args: b.args}, nil
}
