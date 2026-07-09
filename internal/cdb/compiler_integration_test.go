//go:build integration

package cdb

import (
	"bytes"
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

		err := compiler.Compile([]string{"cpx", "run", "./testdata/succeed_compiler.go", "hello"})
		require.NoError(t, err)
		assert.Equal(t, "hello\n", stdout.String())
	})

	t.Run("executes failing command", func(t *testing.T) {
		t.Parallel()

		compiler := &ExecCompiler{
			Bin: "go",
		}

		err := compiler.Compile([]string{"cpx", "run", "./testdata/fail_compiler.go"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "compilation failed")
	})

	t.Run("redirects Stdin", func(t *testing.T) {
		t.Parallel()

		var stdout bytes.Buffer
		stdin := bytes.NewBufferString("hello from stdin")
		compiler := &ExecCompiler{
			Bin:    "go",
			Stdin:  stdin,
			Stdout: &stdout,
		}

		err := compiler.Compile([]string{"cpx", "run", "./testdata/read_stdin_compiler.go"})
		require.NoError(t, err)
		assert.Equal(t, "hello from stdin", stdout.String())
	})

	t.Run("redirects Stderr", func(t *testing.T) {
		t.Parallel()

		var stderr bytes.Buffer
		compiler := &ExecCompiler{
			Bin:    "go",
			Stderr: &stderr,
		}

		err := compiler.Compile([]string{"cpx", "run", "./testdata/write_stderr_compiler.go"})
		require.NoError(t, err)
		assert.Equal(t, "compiler diagnostic message\n", stderr.String())
	})
}
