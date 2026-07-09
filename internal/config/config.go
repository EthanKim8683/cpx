// Package config provides shared configuration for cpx.
package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config holds shared configuration for cpx.
type Config struct {
	// GCC is the absolute path to the GCC compiler driver (required by the gcc shim).
	GCC string `env:"GCC"`
	// GXX is the absolute path to the G++ compiler driver (required by the g++ shim).
	GXX string `env:"GXX"`
	// Clang is the absolute path to the Clang compiler driver (required by the clang shim).
	Clang string `env:"CLANG"`
	// ClangXX is the absolute path to the Clang++ compiler driver (required by the clang++ shim).
	ClangXX string `env:"CLANGXX"`
	// CDB is the destination path for the compilation database JSON file (env: CPX_CDB).
	CDB string `env:"CPX_CDB" envDefault:"./.cpx/cdb.json"`
}

// Load reads configuration from the environment.
// It returns an error only if environment parsing fails.
func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("loading config: %w", err)
	}
	return cfg, nil
}
