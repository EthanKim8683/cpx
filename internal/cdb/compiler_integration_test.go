//go:build integration

package cdb

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecCompiler_Integration(t *testing.T) {
	t.Parallel()

	t.Run("executes successful command", func(t *testing.T) {
		t.Parallel()

		var stdout bytes.Buffer
		compiler := &ExecCompiler{
			Bin:    "go",
			Stdout: &stdout,
		}

		err := compiler.Compile([]string{"cpx", "run", "./testdata/succeed.go"})
		require.NoError(t, err)
	})

	t.Run("executes failing command", func(t *testing.T) {
		t.Parallel()

		compiler := &ExecCompiler{
			Bin: "go",
		}

		err := compiler.Compile([]string{"cpx", "run", "./testdata/fail.go"})
		require.Error(t, err)
		var exitErr *exec.ExitError
		require.ErrorAs(t, err, &exitErr)
		assert.Equal(t, 1, exitErr.ExitCode())
	})

	t.Run("redirects Stdin", func(t *testing.T) {
		t.Parallel()

		var stdout bytes.Buffer
		stdin := bytes.NewBufferString("Hello, world!")
		compiler := &ExecCompiler{
			Bin:    "go",
			Stdin:  stdin,
			Stdout: &stdout,
		}

		err := compiler.Compile([]string{"cpx", "run", "./testdata/stdin.go"})
		require.NoError(t, err)
		assert.Equal(t, "Hello, world!", stdout.String())
	})

	t.Run("redirects Stderr", func(t *testing.T) {
		t.Parallel()

		var stderr bytes.Buffer
		compiler := &ExecCompiler{
			Bin:    "go",
			Stderr: &stderr,
		}

		err := compiler.Compile([]string{"cpx", "run", "./testdata/stderr.go"})
		require.NoError(t, err)
		assert.Equal(t, "Hello, world!\n", stderr.String())
	})
}
