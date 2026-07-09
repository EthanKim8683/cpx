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
			Bin:    "echo",
			Stdout: &stdout,
		}

		err := compiler.Compile([]string{"echo", "hello"})
		require.NoError(t, err)
		assert.Equal(t, "hello\n", stdout.String())
	})

	t.Run("executes failing command", func(t *testing.T) {
		t.Parallel()

		compiler := &ExecCompiler{
			Bin: "false",
		}

		err := compiler.Compile([]string{"false"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "compilation failed")
	})

	t.Run("redirects Stdin", func(t *testing.T) {
		t.Parallel()

		var stdout bytes.Buffer
		stdin := bytes.NewBufferString("hello from stdin")
		compiler := &ExecCompiler{
			Bin:    "cat",
			Stdin:  stdin,
			Stdout: &stdout,
		}

		err := compiler.Compile([]string{"cat"})
		require.NoError(t, err)
		assert.Equal(t, "hello from stdin", stdout.String())
	})

	t.Run("redirects Stderr", func(t *testing.T) {
		t.Parallel()

		var stderr bytes.Buffer
		compiler := &ExecCompiler{
			Bin:    "sh",
			Stderr: &stderr,
		}

		err := compiler.Compile([]string{"sh", "-c", "echo err_msg >&2"})
		require.NoError(t, err)
		assert.Equal(t, "err_msg\n", stderr.String())
	})
}
