package clang_test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/EthanKim8683/cpx/internal/bundler/clang"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func TestBundler(t *testing.T) {
	t.Parallel()

	g := goldie.New(t)

	flags := []string{
		"-std=c++17",
		"-I./testdata/include",
		"-o",
		"./testdata/src/main",
	}
	b := clang.NewBundler("clang", flags)

	bundle, err := b.Bundle(t.Context(), "./testdata/src/main.cpp")
	require.NoError(t, err)
	g.Assert(t, t.Name(), []byte(bundle))

	var (
		stdin  = bytes.NewBuffer([]byte(bundle))
		stderr = bytes.NewBuffer([]byte{})
	)
	cmd := exec.CommandContext(
		t.Context(),
		"clang",
		append(flags, "-o", "/dev/null", "-x", "c++", "-")...,
	)
	cmd.Stdin = stdin
	cmd.Stderr = stderr
	require.NoError(t, cmd.Run(), stderr.String())
}
