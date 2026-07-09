// Package config provides shared configuration for cpx.
package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config holds shared configuration for cpx.
type Config struct {
	// GCC is the absolute path to the GCC compiler driver (required).
	GCC string `env:"GCC"`
	// GXX is the absolute path to the G++ compiler driver (optional).
	GXX string `env:"GXX"`
	// Clang is the absolute path to the Clang compiler driver (required).
	Clang string `env:"CLANG"`
	// ClangXX is the absolute path to the Clang++ compiler driver (optional).
	ClangXX string `env:"CLANGXX"`
	// CDB is the destination path to write the compilation database JSON file.
	CDB string `env:"CPX_CDB" envDefault:"./.cpx/cdb.json"`
}

// Load reads configuration from the environment.
// It returns an error if any required value is missing or empty.
func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("loading config: %w", err)
	}
	return cfg, nil
}
