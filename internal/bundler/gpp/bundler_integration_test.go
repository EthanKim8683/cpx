// //go:build integration

package gpp_test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/EthanKim8683/cpx/internal/bundler/gpp"
	"github.com/EthanKim8683/cpx/internal/config"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBundler(t *testing.T) {
	t.Parallel()

	cfg := config.Load()
	var (
		executable = cfg.Gpp
		flags      = []string{
			"-std=c++17",
			"-I./testdata/include",
			"-o",
			"./testdata/main",
		}
	)
	if executable == "" {
		t.Skip("CPX_GPP is not set")
	}

	g := goldie.New(t)

	t.Run("happy path", func(t *testing.T) {
		b, err := gpp.NewBundler(cfg, append([]string{executable, "./testdata/happy_path.cpp"}, flags...))
		require.NoError(t, err)
		bundle, err := b.Bundle(t.Context())
		require.NoError(t, err)
		g.Assert(t, t.Name(), []byte(bundle))

		stdin := bytes.NewBuffer([]byte(bundle))
		var stderr bytes.Buffer
		cmd := exec.CommandContext(
			t.Context(),
			executable,
			append(flags, "-o", "/dev/null", "-x", "c++", "-")...,
		)
		cmd.Stdin = stdin
		cmd.Stderr = &stderr
		require.NoError(t, cmd.Run(), stderr.String())
	})

	t.Run("no arguments", func(t *testing.T) {
		_, err := gpp.NewBundler(cfg, nil)
		assert.ErrorContains(t, err, "no arguments provided")
	})
}
