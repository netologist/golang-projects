package main

import "sync"

// Job is a unit of work.
type Job struct {
	ID    int
	Value int
}

// Result is the outcome of a Job.
type Result struct {
	JobID  int
	Output int
}

// Run starts `workers` goroutines consuming jobs and returns a results channel.
func Run(workers int, jobs <-chan Job) <-chan Result {
	results := make(chan Result, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				results <- Result{JobID: j.ID, Output: j.Value * j.Value}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}
