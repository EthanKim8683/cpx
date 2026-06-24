//go:build integration

package gpp_test

import (
	"bytes"
	"os/exec"
	"slices"
	"testing"

	"github.com/EthanKim8683/cpx/internal/bundler/gpp"
	"github.com/EthanKim8683/cpx/internal/config"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

var compileFlags = []string{
	"-std=c++17",
	"-o",
	"./testdata/main",
}

func TestBundler(t *testing.T) {
	t.Parallel()

	cfg := config.Load()
	if cfg.Gpp == "" {
		t.Skip("CPX_GPP is not set")
	}
	if cfg.Clangpp == "" {
		t.Skip("CPX_CLANGPP is not set")
	}

	g := goldie.New(t)

	t.Run("happy path", func(t *testing.T) {
		var bundle string
		{
			compileArgs := slices.Concat(
				[]string{cfg.Gpp, "./testdata/happy_path.cpp"},
				compileFlags,
			)
			b, err := gpp.NewBundler(cfg, compileArgs)
			require.NoError(t, err)
			bundle, err = b.Bundle(t.Context())
			require.NoError(t, err)
		}
		g.Assert(t, t.Name(), []byte(bundle))

		{
			compileArgs := slices.Concat(
				compileFlags,
				[]string{"-o", "/dev/null", "-x", "c++", "-"},
			)
			stdin := bytes.NewBufferString(bundle)
			var stderr bytes.Buffer
			cmd := exec.CommandContext(t.Context(), cfg.Gpp, compileArgs...)
			cmd.Stdin = stdin
			cmd.Stderr = &stderr
			require.NoError(t, cmd.Run(), stderr.String())
		}
	})
}
