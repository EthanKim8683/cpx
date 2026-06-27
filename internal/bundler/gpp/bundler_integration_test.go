//go:build integration

package gpp_test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/EthanKim8683/cpx/internal/bundler/gpp"
	"github.com/EthanKim8683/cpx/internal/config"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

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
		b := gpp.NewBundler(
			cfg.Clangpp,
			[]string{
				"-I./testdata/include",
				"-o./testdata/happy_path",
				"./testdata/happy_path.cpp",
			},
		)
		bundle, err := b.Bundle(t.Context())
		require.NoError(t, err)
		g.Assert(t, t.Name(), []byte(bundle))

		stdin := bytes.NewBufferString(bundle)
		var stderr bytes.Buffer
		cmd := exec.CommandContext(
			t.Context(),
			cfg.Gpp,
			"-xc++",
			"-o/dev/null",
			"-",
		)
		cmd.Stdin = stdin
		cmd.Stderr = &stderr
		require.NoError(t, cmd.Run(), stderr.String())
	})
}
