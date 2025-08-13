package worker

import (
	"context"

	"mastercli-starter/internal/job"
)

type Result struct {
	JobID string
	Err   error
}

type Worker struct {
	id int
}

func New(id int) *Worker { return &Worker{id: id} }

func (w *Worker) Run(ctx context.Context, jobs <-chan job.Job, results chan<- Result) {
	for {
		select {
		case <-ctx.Done():
			return
		case j, ok := <-jobs:
			if !ok {
				return
			}
			err := j.Do(ctx)
			results <- Result{JobID: j.ID(), Err: err}
		}
	}
}
