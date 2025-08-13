package mastercli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/mastercli/internal/config"
	"example.com/mastercli/internal/logger"
	"github.com/spf13/cobra"
)

var (
	cfgPath string
	rootCmd = &cobra.Command{
		Use:   "mastercli",
		Short: "Master framework + CLI in Go",
		Long:  "A starter master-worker framework with a modern CLI (Cobra/Viper).",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "Path to config file (default: ./configs/config.yaml)")
}

func Execute() {
	// Load config before running commands
	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}
	log := logger.Init(cfg.App.LogLevel)
	log.Info().Str("app", cfg.App.Name).Msg("starting")

	// Graceful shutdown wiring for any subcommand that cares
	ctx, stop := signal.NotifyContext(os.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Inject context via command context
	rootCmd.SetContext(ctx)

	// Register subcommands
	rootCmd.AddCommand(startCmd(cfg))
	rootCmd.AddCommand(runCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("command failed")
		os.Exit(1)
	}
	log.Info().Dur("uptime", time.Since(time.Now())).Msg("exiting")
}
