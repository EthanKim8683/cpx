//go:build integration

package clangpp_test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/EthanKim8683/cpx/internal/bundler/clangpp"
	"github.com/EthanKim8683/cpx/internal/config"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func TestBundler(t *testing.T) {
	t.Parallel()

	cfg := config.Load()
	var (
		executable = cfg.Clangpp
		flags      = []string{
			"-std=c++17",
			"-I./testdata/include",
			"-o",
			"./testdata/src/main",
		}
	)
	if executable == "" {
		t.Skip("CPX_CLANGPP is not set")
	}

	g := goldie.New(t)

	b, err := clangpp.NewBundler(append([]string{executable, "./testdata/src/main.cpp"}, flags...))
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
}
