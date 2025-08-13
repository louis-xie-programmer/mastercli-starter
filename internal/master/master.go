package master

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"example.com/mastercli/internal/job"
	"example.com/mastercli/internal/worker"
	"example.com/mastercli/internal/logger"
)

type Manager struct {
	workers    int
	queueSize  int
	maxRetries int
	backoff    time.Duration

	jobs    chan job.Job
	results chan worker.Result
}

func NewManager(workers, queueSize, maxRetries int, backoff time.Duration) *Manager {
	return &Manager{
		workers:    workers,
		queueSize:  queueSize,
		maxRetries: maxRetries,
		backoff:    backoff,
		jobs:       make(chan job.Job, queueSize),
		results:    make(chan worker.Result, queueSize),
	}
}

// Submit enqueues a job (non-blocking if buffer available).
func (m *Manager) Submit(j job.Job) error {
	select {
	case m.jobs <- j:
		return nil
	default:
		return errors.New("job queue full")
	}
}

func (m *Manager) Jobs() chan<- job.Job { return m.jobs }

func (m *Manager) Start(ctx context.Context) {
	log := logger.L()
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < m.workers; i++ {
		wg.Add(1)
		w := worker.New(i + 1)
		go func() {
			defer wg.Done()
			w.Run(ctx, m.jobs, m.results)
		}()
	}

	// Collector
	go func() {
		wg.Wait()
		close(m.results)
	}()

	// Retry loop & logging
	go func() {
		for r := range m.results {
			if r.Err == nil {
				log.Info().Str("job_id", r.JobID).Msg("completed")
				continue
			}
			log.Warn().Str("job_id", r.JobID).Err(r.Err).Msg("failed")
			// naive retry by re-enqueueing with exponential backoff
			for attempt := 1; attempt <= m.maxRetries; attempt++ {
				backoff := m.backoff * time.Duration(math.Pow(2, float64(attempt-1)))
				select {
				case <-ctx.Done():
					return
				case <-time.After(backoff):
					log.Debug().Str("job_id", r.JobID).Int("attempt", attempt).Dur("backoff", backoff).Msg("retrying")
					// In a real impl we would clone the job; for demo the job instance retries itself
					// by being re-submitted from the worker loop. Here we just log and mark as success.
					// (Alternatively keep a map[ID]job.Job to requeue the same job instance.)
				}
			}
			log.Error().Str("job_id", r.JobID).Msg("exhausted retries")
		}
	}()
}

func (m *Manager) Stop() {
	close(m.jobs)
}

func GenerateDemoJobs(n int) []job.Job {
	jobs := make([]job.Job, 0, n)
	for i := 0; i < n; i++ {
		failOnce := i%5 == 0
		jobs = append(jobs, &job.SimpleJob{
			JobID:    fmt.Sprintf("job-%03d", i+1),
			Payload:  fmt.Sprintf("payload-%d", i+1),
			Duration: time.Duration(100+10*i) * time.Millisecond,
			FailOnce: failOnce,
		})
	}
	return jobs
}
