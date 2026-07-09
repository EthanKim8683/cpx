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
}
