package config

import (
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg, err := Load(env.Options{
			Environment: map[string]string{
				"GCC":     "/usr/bin/gcc",
				"GXX":     "/usr/bin/g++",
				"CLANG":   "/usr/bin/clang",
				"CLANGXX": "/usr/bin/clang++",
			},
		})
		require.NoError(t, err)
		assert.Equal(t, "/usr/bin/gcc", cfg.GCC)
		assert.Equal(t, "/usr/bin/g++", cfg.GXX)
		assert.Equal(t, "/usr/bin/clang", cfg.Clang)
		assert.Equal(t, "/usr/bin/clang++", cfg.ClangXX)
		assert.Equal(t, "./cpx/cdb.json", cfg.CDB)
	})

	t.Run("missing GCC", func(t *testing.T) {
		_, err := Load(env.Options{
			Environment: map[string]string{
				"GXX":     "/usr/bin/g++",
				"CLANG":   "/usr/bin/clang",
				"CLANGXX": "/usr/bin/clang++",
			},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "loading config")
	})

	t.Run("missing GXX (allowed)", func(t *testing.T) {
		cfg, err := Load(env.Options{
			Environment: map[string]string{
				"GCC":     "/usr/bin/gcc",
				"CLANG":   "/usr/bin/clang",
				"CLANGXX": "/usr/bin/clang++",
			},
		})
		require.NoError(t, err)
		assert.Empty(t, cfg.GXX)
	})

	t.Run("missing CLANG", func(t *testing.T) {
		_, err := Load(env.Options{
			Environment: map[string]string{
				"GCC":     "/usr/bin/gcc",
				"GXX":     "/usr/bin/g++",
				"CLANGXX": "/usr/bin/clang++",
			},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "loading config")
	})

	t.Run("missing CLANGXX (allowed)", func(t *testing.T) {
		cfg, err := Load(env.Options{
			Environment: map[string]string{
				"GCC":   "/usr/bin/gcc",
				"GXX":   "/usr/bin/g++",
				"CLANG": "/usr/bin/clang",
			},
		})
		require.NoError(t, err)
		assert.Empty(t, cfg.ClangXX)
	})

	t.Run("empty GCC", func(t *testing.T) {
		_, err := Load(env.Options{
			Environment: map[string]string{
				"GCC":     "",
				"GXX":     "/usr/bin/g++",
				"CLANG":   "/usr/bin/clang",
				"CLANGXX": "/Users/bin/clang++",
			},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "loading config")
	})

	t.Run("all missing", func(t *testing.T) {
		_, err := Load(env.Options{
			Environment: map[string]string{},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "loading config")
	})
}
