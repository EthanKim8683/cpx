package gcc_test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/EthanKim8683/cpx/internal/bundler/gcc"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func TestBundler(t *testing.T) {
	t.Parallel()

	g := goldie.New(t)

	var (
		command = "g++"
		flags   = []string{
			"-std=c++17",
			"-I./testdata/include",
			"-o",
			"./testdata/src/main",
		}
		args = append(flags, "./testdata/src/main.cpp")
	)
	b := gcc.NewBundler(append([]string{command}, args...))

	bundle, err := b.Bundle(t.Context())
	require.NoError(t, err)
	g.Assert(t, t.Name(), []byte(bundle))

	var (
		stdin  = bytes.NewBuffer([]byte(bundle))
		stderr = bytes.NewBuffer([]byte{})
	)
	cmd := exec.CommandContext(
		t.Context(),
		command,
		append(flags, "-o", "/dev/null", "-x", "c++", "-")...,
	)
	cmd.Stdin = stdin
	cmd.Stderr = stderr
	require.NoError(t, cmd.Run(), stderr.String())
}
