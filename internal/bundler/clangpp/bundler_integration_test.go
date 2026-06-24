//go:build integration

package clangpp_test

import (
	"bytes"
	"os/exec"
	"slices"
	"testing"

	"github.com/EthanKim8683/cpx/internal/bundler/clangpp"
	"github.com/EthanKim8683/cpx/internal/config"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var compileFlags = []string{
	"-std=c++17",
	"-I./testdata/include",
	"-o",
	"./testdata/src/main",
}

func TestBundler(t *testing.T) {
	t.Parallel()

	cfg := config.Load()
	if cfg.Clangpp == "" {
		t.Skip("CPX_CLANGPP is not set")
	}

	g := goldie.New(t)

	t.Run("happy path", func(t *testing.T) {
		var bundle string
		{
			compileArgs := slices.Concat(
				[]string{cfg.Clangpp, "./testdata/src/happy_path.cpp"},
				compileFlags,
			)
			b, err := clangpp.NewBundler(compileArgs)
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
			cmd := exec.CommandContext(t.Context(), cfg.Clangpp, compileArgs...)
			cmd.Stdin = stdin
			cmd.Stderr = &stderr
			require.NoError(t, cmd.Run(), stderr.String())
		}
	})

	t.Run("g++ only", func(t *testing.T) {
		compileArgs := slices.Concat(
			[]string{cfg.Clangpp, "./testdata/src/g++_only.cpp"},
			compileFlags,
		)
		b, err := clangpp.NewBundler(compileArgs)
		require.NoError(t, err)
		_, err = b.Bundle(t.Context())
		assert.Error(t, err)
	})
}
