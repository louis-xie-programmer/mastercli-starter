package mastercli

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"example.com/mastercli/internal/logger"
	"example.com/mastercli/internal/master"
	"example.com/mastercli/internal/job"
	"github.com/spf13/cobra"
)

type startOptions struct {
	fromFile string
	jobs     int
}

func startCmd(appCfg *struct{ App struct{ Name string `yaml:"name"`; LogLevel string `yaml:"log_level"` `yaml:"log_level"`}; Master struct{ Workers int `yaml:"workers"`; QueueSize int `yaml:"queue_size"`; MaxRetries int `yaml:"max_retries"`; BackoffMS int `yaml:"backoff_ms"` } `yaml:"master"` }) *cobra.Command {
	opts := &startOptions{}
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the master and process jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			log := logger.L()

			mgr := master.NewManager(appCfg.Master.Workers, appCfg.Master.QueueSize, appCfg.Master.MaxRetries, time.Duration(appCfg.Master.BackoffMS)*time.Millisecond)
			mgr.Start(ctx)

			// Feed jobs
			go func() {
				defer mgr.Stop()
				if opts.fromFile != "" {
					f, err := os.Open(opts.fromFile)
					if err != nil {
						log.Error().Err(err).Msg("open file")
						return
					}
					defer f.Close()
					scanner := bufio.NewScanner(f)
					i := 0
					for scanner.Scan() {
						i++
						payload := scanner.Text()
						j := &job.SimpleJob{JobID: fmt.Sprintf("file-%03d", i), Payload: payload, Duration: 200*time.Millisecond}
						if err := mgr.Submit(j); err != nil {
							log.Warn().Err(err).Msg("submit")
						}
					}
					return
				}
				for _, j := range master.GenerateDemoJobs(opts.jobs) {
					if err := mgr.Submit(j); err != nil { log.Warn().Err(err).Msg("submit") }
				}
			}()

			// Block until context cancelled
			<-ctx.Done()
			log.Info().Msg("shutdown signal received")
			return nil
		},
	}
	cmd.Flags().StringVarP(&opts.fromFile, "from-file", "f", "", "Read job payloads from a text file (one per line)")
	cmd.Flags().IntVarP(&opts.jobs, "jobs", "n", 20, "Number of demo jobs to generate if no file is provided")
	return cmd
}
