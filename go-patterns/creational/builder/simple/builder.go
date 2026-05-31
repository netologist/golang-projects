package main

import "fmt"

// Query is an immutable SQL query.
type Query struct {
	table string
	cols  []string
	where string
}

func (q Query) String() string {
	cols := "*"
	if len(q.cols) > 0 {
		cols = fmt.Sprintf("%v", q.cols)
	}
	s := fmt.Sprintf("SELECT %s FROM %s", cols, q.table)
	if q.where != "" {
		s += " WHERE " + q.where
	}
	return s
}

// QueryBuilder builds a Query.
type QueryBuilder struct {
	table string
	cols  []string
	where string
}

// NewQueryBuilder starts a builder for the given table.
func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{table: table}
}

// Select sets the columns to select.
func (b *QueryBuilder) Select(cols ...string) *QueryBuilder {
	b.cols = cols
	return b
}

// Where sets the WHERE clause.
func (b *QueryBuilder) Where(cond string) *QueryBuilder {
	b.where = cond
	return b
}

// Build validates and produces the immutable Query.
func (b *QueryBuilder) Build() (Query, error) {
	if b.table == "" {
		return Query{}, fmt.Errorf("table is required")
	}
	return Query{table: b.table, cols: b.cols, where: b.where}, nil
}
