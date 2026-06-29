// Package config loads shared tool paths from environment variables at CLI startup.
//
// Direnv loads .env via .envrc. This package does not parse .env files.
package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config holds compiler tool paths read from environment variables.
type Config struct {
	// GCC is the path to the GCC driver.
	GCC string `env:"GCC,notEmpty"`

	// Clang is the path to the Clang driver.
	Clang string `env:"CLANG,notEmpty"`

	// ClangTblgen is the path to the clang-tblgen binary.
	ClangTblgen string `env:"CLANG_TBLGEN,notEmpty"`
}

// Load reads required tool paths from the environment.
// It returns an error if any variable is missing or empty.
func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("config: %w", err)
	}
	return cfg, nil
}
