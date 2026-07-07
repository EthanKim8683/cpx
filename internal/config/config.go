// Package config provides shared configuration for cpx.
package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config holds shared configuration for cpx.
type Config struct {
	// GCC is the path to the GCC driver.
	GCC string `env:"GCC,notEmpty"`
	// Clang is the path to the Clang driver.
	Clang string `env:"CLANG,notEmpty"`
}

// Load reads configuration from the environment.
// It returns an error if any required value is missing or empty.
// Pass env.Options to override the environment (e.g. for testing).
func Load(opts ...env.Options) (Config, error) {
	var cfg Config
	var err error
	if len(opts) > 0 {
		cfg, err = env.ParseAsWithOptions[Config](opts[0])
	} else {
		cfg, err = env.ParseAs[Config]()
	}
	if err != nil {
		return Config{}, fmt.Errorf("loading config: %w", err)
	}
	return cfg, nil
}
