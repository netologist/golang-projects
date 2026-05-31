package main

import "sync"

// DB simulates an expensive database connection pool.
type DB struct{ dsn string }

// Ping reports whether the DB is reachable.
func (d *DB) Ping() bool { return true }

var (
	dbInstance *DB
	dbOnce     sync.Once
)

// GetDB returns the singleton DB instance.
func GetDB() *DB {
	dbOnce.Do(func() {
		dbInstance = &DB{dsn: "postgres://localhost/app"}
	})
	return dbInstance
}

// Reset tears down the singleton for test isolation.
func Reset() {
	dbOnce = sync.Once{}
	dbInstance = nil
}
