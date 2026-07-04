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
	// GCCBaseURL is the base URL for fetching GCC sources over HTTP.
	GCCBaseURL string `env:"GCC_BASE_URL,notEmpty"`
	// Clang is the path to the Clang driver.
	Clang string `env:"CLANG,notEmpty"`
	// ClangTblgen is the path to the clang-tblgen binary.
	ClangTblgen string `env:"CLANG_TBLGEN,notEmpty"`
	// LLVMBaseURL is the base URL for fetching LLVM sources over HTTP.
	LLVMBaseURL string `env:"LLVM_BASE_URL,notEmpty"`
}

// Load reads configuration from the environment.
// It returns an error if any required value is missing or empty.
func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("config: %w", err)
	}
	return cfg, nil
}
