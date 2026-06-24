//go:build integration

package gpp_test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/EthanKim8683/cpx/internal/bundler/gpp"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func TestBundler(t *testing.T) {
	t.Parallel()

	g := goldie.New(t)

	var (
		executable = "/opt/homebrew/bin/g++-16"
		flags      = []string{
			"-std=c++17",
			"-I./testdata/include",
			"-o",
			"./testdata/main",
		}
		args = append(flags, "./testdata/main.cpp")
	)
	b := gpp.NewBundler(append([]string{executable}, args...))

	bundle, err := b.Bundle(t.Context())
	require.NoError(t, err)
	g.Assert(t, t.Name(), []byte(bundle))

	var (
		stdin  = bytes.NewBuffer([]byte(bundle))
		stderr = bytes.NewBuffer([]byte{})
	)
	cmd := exec.CommandContext(
		t.Context(),
		executable,
		append(flags, "-o", "/dev/null", "-x", "c++", "-")...,
	)
	cmd.Stdin = stdin
	cmd.Stderr = stderr
	require.NoError(t, cmd.Run(), stderr.String())
}
