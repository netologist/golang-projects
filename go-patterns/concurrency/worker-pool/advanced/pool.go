package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// Job is a unit of work with an arbitrary payload.
type Job struct {
	ID      int
	Payload any
}

// Result is the outcome of processing a Job.
type Result struct {
	JobID  int
	Output any
	Err    error
}

// Pool is a bounded worker pool with graceful shutdown.
type Pool struct {
	jobs    chan Job
	results chan Result
	wg      sync.WaitGroup
	active  atomic.Int64
	done    chan struct{}
}

// New starts a pool with the given worker count and queue size.
func New(ctx context.Context, workers, queueSize int) *Pool {
	p := &Pool{
		jobs:    make(chan Job, queueSize),
		results: make(chan Result, queueSize),
		done:    make(chan struct{}),
	}
	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go p.worker(ctx)
	}
	go func() {
		p.wg.Wait()
		close(p.results)
		close(p.done)
	}()
	return p
}

func (p *Pool) worker(ctx context.Context) {
	defer p.wg.Done()
	for {
		select {
		case job, ok := <-p.jobs:
			if !ok {
				return
			}
			p.active.Add(1)
			p.results <- Result{JobID: job.ID, Output: job.Payload}
			p.active.Add(-1)
		case <-ctx.Done():
			return
		}
	}
}

// Submit enqueues a job, honouring context cancellation.
func (p *Pool) Submit(ctx context.Context, job Job) error {
	select {
	case p.jobs <- job:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("submit job %d: %w", job.ID, ctx.Err())
	}
}

// Results returns the channel of completed results.
func (p *Pool) Results() <-chan Result { return p.results }

// ActiveWorkers reports how many workers are currently processing.
func (p *Pool) ActiveWorkers() int64 { return p.active.Load() }

// Shutdown stops accepting jobs and waits for workers to drain.
func (p *Pool) Shutdown() {
	close(p.jobs)
	<-p.done
}
