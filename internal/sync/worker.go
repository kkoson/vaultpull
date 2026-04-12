package sync

import (
	"context"
	"fmt"
	"sync"
)

// WorkerPool manages a pool of goroutines for concurrent profile syncing.
type WorkerPool struct {
	workers   int
	jobs      chan workerJob
	results   chan WorkerResult
	wg        sync.WaitGroup
}

type workerJob struct {
	profileName string
	runFn       func(ctx context.Context) error
}

// WorkerResult holds the outcome of a single worker job.
type WorkerResult struct {
	ProfileName string
	Err         error
}

// NewWorkerPool creates a WorkerPool with the given concurrency level.
// workers must be >= 1; if not, it defaults to 1.
func NewWorkerPool(workers int) *WorkerPool {
	if workers < 1 {
		workers = 1
	}
	return &WorkerPool{
		workers: workers,
		jobs:    make(chan workerJob),
		results: make(chan WorkerResult),
	}
}

// Start launches worker goroutines and returns a results channel.
func (p *WorkerPool) Start(ctx context.Context) <-chan WorkerResult {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for job := range p.jobs {
				err := job.runFn(ctx)
				p.results <- WorkerResult{ProfileName: job.profileName, Err: err}
			}
		}()
	}
	go func() {
		p.wg.Wait()
		close(p.results)
	}()
	return p.results
}

// Submit enqueues a job for the given profile. Returns an error if ctx is done.
func (p *WorkerPool) Submit(ctx context.Context, profileName string, fn func(ctx context.Context) error) error {
	select {
	case p.jobs <- workerJob{profileName: profileName, runFn: fn}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("worker pool: context cancelled while submitting %q", profileName)
	}
}

// Close signals that no more jobs will be submitted.
func (p *WorkerPool) Close() {
	close(p.jobs)
}
