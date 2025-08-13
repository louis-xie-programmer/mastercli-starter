package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	App struct {
		Name     string `yaml:"name"`
		LogLevel string `yaml:"log_level"`
	} `yaml:"app"`
	Master struct {
		Workers    int `yaml:"workers"`
		QueueSize  int `yaml:"queue_size"`
		MaxRetries int `yaml:"max_retries"`
		BackoffMS  int `yaml:"backoff_ms"`
	} `yaml:"master"`
}

func Load(configPath string) (*AppConfig, error) {
	v := viper.New()
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.AddConfigPath("./configs")
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}
	v.SetEnvPrefix("MASTERCLI")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &cfg, nil
}
