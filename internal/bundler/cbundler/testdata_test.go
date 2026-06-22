package cbundler

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func testdataRoot(t *testing.T) string {
	t.Helper()

	absPath, err := filepath.Abs(filepath.Join("testdata"))
	require.NoError(t, err)
	return absPath
}
