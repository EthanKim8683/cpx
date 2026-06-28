//go:build integration

package clang_test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/EthanKim8683/cpx/internal/bundler/clang"
	"github.com/EthanKim8683/cpx/internal/config"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBundler(t *testing.T) {
	t.Parallel()

	cfg := config.Load()
	if cfg.Clangpp == "" {
		t.Skip("CPX_CLANGPP is not set")
	}

	g := goldie.New(t)

	t.Run("happy path", func(t *testing.T) {
		b := clang.NewBundler(
			cfg.Clangpp,
			[]string{
				"-I./testdata/include",
				"-o./testdata/src/happy_path",
				"./testdata/src/happy_path.cpp",
			},
		)
		bundle, err := b.Bundle(t.Context())
		require.NoError(t, err)
		g.Assert(t, t.Name(), []byte(bundle))

		stdin := bytes.NewBufferString(bundle)
		var stderr bytes.Buffer
		cmd := exec.CommandContext(
			t.Context(),
			cfg.Clangpp,
			"-xc++",
			"-o/dev/null",
			"-",
		)
		cmd.Stdin = stdin
		cmd.Stderr = &stderr
		require.NoError(t, cmd.Run(), stderr.String())
	})

	t.Run("g++ only", func(t *testing.T) {
		b := clang.NewBundler(
			cfg.Clangpp,
			[]string{
				"./testdata/src/g++_only.cpp",
				"-o./testdata/src/g++_only",
			},
		)
		_, err := b.Bundle(t.Context())
		assert.ErrorContains(t, err, "bits/stdc++.h")
	})

	t.Run("multiple files", func(t *testing.T) {
		b := clang.NewBundler(
			cfg.Clangpp,
			[]string{
				"./testdata/src/multiple_files.cpp",
				"./testdata/src/multiple_files.cpp",
			},
		)
		_, err := b.Bundle(t.Context())
		assert.NoError(t, err)
	})
}
