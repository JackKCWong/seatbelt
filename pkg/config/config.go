package config

import "github.com/spf13/viper"

type Config struct {
	Action   string   `mapstructure:"action"`
	Patterns []string `mapstructure:"patterns"`
}

func Load() (Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}