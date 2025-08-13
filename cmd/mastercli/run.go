package mastercli

import (
	"context"
	"fmt"
	"time"

	"example.com/mastercli/internal/job"
	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	var (
		payload string
		durMs   int
		fail    bool
	)
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a single job synchronously and print the result",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), time.Duration(durMs+1000)*time.Millisecond)
			defer cancel()
			j := &job.SimpleJob{JobID: "single", Payload: payload, Duration: time.Duration(durMs) * time.Millisecond, FailOnce: fail}
			fmt.Printf("Running %s...\n", j.String())
			if err := j.Do(ctx); err != nil {
				return err
			}
			fmt.Println("OK")
			return nil
		},
	}
	cmd.Flags().StringVarP(&payload, "payload", "p", "hello", "Payload string")
	cmd.Flags().IntVarP(&durMs, "duration", "d", 500, "Duration in milliseconds")
	cmd.Flags().BoolVarP(&fail, "fail-once", "F", false, "Simulate a transient failure once")
	return cmd
}
