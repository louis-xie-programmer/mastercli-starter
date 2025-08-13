package mastercli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"mastercli-starter/internal/config"
	"mastercli-starter/internal/logger"
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
	// 加载配置
	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}
	log := logger.Init(cfg.App.LogLevel)
	log.Info().Str("app", cfg.App.Name).Msg("starting")

	//  创建一个cancelable context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 设置context到rootCmd
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
