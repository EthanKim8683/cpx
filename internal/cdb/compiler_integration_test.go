//go:build integration

package cdb

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// compileHelper compiles the given Go source file to a temporary executable for testing.
func compileHelper(t *testing.T, src string) string {
	t.Helper()
	tmpDir := t.TempDir()
	binName := filepath.Base(src)
	binName = binName[:len(binName)-len(filepath.Ext(binName))]
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	dest := filepath.Join(tmpDir, binName)
	//nolint:gosec // executing local test source compiled dynamically in tests
	cmd := exec.Command("go", "build", "-o", dest, src)
	err := cmd.Run()
	require.NoError(t, err)
	return dest
}

func TestExecCompiler_Integration(t *testing.T) {
	t.Parallel()

	t.Run("executes successful command", func(t *testing.T) {
		t.Parallel()

		helperBin := compileHelper(t, "./testdata/echo.go")

		var stdout bytes.Buffer
		compiler := &ExecCompiler{
			Bin:    helperBin,
			Stdout: &stdout,
		}

		err := compiler.Compile([]string{helperBin, "hello"})
		require.NoError(t, err)
		assert.Equal(t, "hello\n", stdout.String())
	})

	t.Run("executes failing command", func(t *testing.T) {
		t.Parallel()

		helperBin := compileHelper(t, "./testdata/false.go")

		compiler := &ExecCompiler{
			Bin: helperBin,
		}

		err := compiler.Compile([]string{helperBin})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "compilation failed")
	})

	t.Run("redirects Stdin", func(t *testing.T) {
		t.Parallel()

		helperBin := compileHelper(t, "./testdata/cat.go")

		var stdout bytes.Buffer
		stdin := bytes.NewBufferString("hello from stdin")
		compiler := &ExecCompiler{
			Bin:    helperBin,
			Stdin:  stdin,
			Stdout: &stdout,
		}

		err := compiler.Compile([]string{helperBin})
		require.NoError(t, err)
		assert.Equal(t, "hello from stdin", stdout.String())
	})

	t.Run("redirects Stderr", func(t *testing.T) {
		t.Parallel()

		helperBin := compileHelper(t, "./testdata/stderr.go")

		var stderr bytes.Buffer
		compiler := &ExecCompiler{
			Bin:    helperBin,
			Stderr: &stderr,
		}

		err := compiler.Compile([]string{helperBin})
		require.NoError(t, err)
		assert.Equal(t, "err_msg\n", stderr.String())
	})
}
