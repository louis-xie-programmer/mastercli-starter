package job

import (
	"context"
	"fmt"
	"time"
)

// Job represents a unit of work.
type Job interface {
	ID() string
	Do(ctx context.Context) error
	String() string
}

// SimpleJob is an example job that sleeps for Duration and optionally fails once.
type SimpleJob struct {
	JobID    string
	Payload  string
	Duration time.Duration
	FailOnce bool
	failed   bool
}

func (j *SimpleJob) ID() string { return j.JobID }

func (j *SimpleJob) String() string {
	return fmt.Sprintf("SimpleJob{id=%s,payload=%s,dur=%s}", j.JobID, j.Payload, j.Duration)
}

func (j *SimpleJob) Do(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(j.Duration):
		// continue
	}
	if j.FailOnce && !j.failed {
		j.failed = true
		return fmt.Errorf("transient failure for job %s", j.JobID)
	}
	return nil
}
