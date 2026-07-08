// Package config provides shared configuration for cpx.
package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config holds shared configuration for cpx.
type Config struct {
	// GCC is the path to the gcc driver.
	GCC string `env:"GCC,notEmpty"`
	// GXX is the path to the g++ driver.
	GXX string `env:"GXX"`
	// Clang is the path to the clang driver.
	Clang string `env:"CLANG,notEmpty"`
	// ClangXX is the path to the clang++ driver.
	ClangXX string `env:"CLANGXX"`
	// CDB is the path to the compilation database file.
	CDB string `env:"CPX_CDB" envDefault:"./cpx/cdb.json"`
}

// Load reads configuration from the environment.
// It returns an error if any required value is missing or empty.
// Pass at most one env.Options to override the environment (e.g. for testing).
func Load(opts ...env.Options) (Config, error) {
	if len(opts) > 1 {
		return Config{}, fmt.Errorf("loading config: too many options arguments (expected at most 1, got %d)", len(opts))
	}

	var cfg Config
	var err error
	if len(opts) == 1 {
		cfg, err = env.ParseAsWithOptions[Config](opts[0])
	} else {
		cfg, err = env.ParseAs[Config]()
	}
	if err != nil {
		return Config{}, fmt.Errorf("loading config: %w", err)
	}
	return cfg, nil
}
