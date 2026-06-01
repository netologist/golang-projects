package main

import (
	"errors"
	"fmt"
	"sync"
)

// ErrClosed is returned when operations are attempted on a closed WAL.
var ErrClosed = errors.New("wal: closed")

// Entry is a single WAL record.
type Entry struct {
	LSN       uint64
	Operation string // e.g., "SET", "DEL"
	Key       string
	Value     string
}

// ApplyFn applies a WAL entry to the state machine.
type ApplyFn func(e Entry)

// WAL is a write-ahead log backed by an in-memory buffer.
type WAL struct {
	mu            sync.Mutex
	entries       []Entry
	seq           uint64
	checkpointLSN uint64
	applyFn       ApplyFn
	closed        bool
}

// New creates a WAL with the given apply function.
func New(fn ApplyFn) *WAL {
	return &WAL{applyFn: fn}
}

// Append writes an entry to the log and applies it.
func (w *WAL) Append(op, key, value string) (uint64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return 0, ErrClosed
	}
	w.seq++
	e := Entry{LSN: w.seq, Operation: op, Key: key, Value: value}
	w.entries = append(w.entries, e)
	w.applyFn(e)
	return w.seq, nil
}

// Checkpoint marks all entries up to and including lsn as checkpointed.
func (w *WAL) Checkpoint(lsn uint64) {
	w.mu.Lock()
	w.checkpointLSN = lsn
	w.mu.Unlock()
}

// Recover replays all entries after the checkpoint LSN, calling applyFn for each.
func (w *WAL) Recover() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, e := range w.entries {
		if e.LSN > w.checkpointLSN {
			w.applyFn(e)
		}
	}
	return nil
}

// Close marks the WAL as closed.
func (w *WAL) Close() {
	w.mu.Lock()
	w.closed = true
	w.mu.Unlock()
}

// Entries returns all log entries (for inspection/testing).
func (w *WAL) Entries() []Entry {
	w.mu.Lock()
	out := append([]Entry{}, w.entries...)
	w.mu.Unlock()
	return out
}

var _ = fmt.Sprintf // keep import
