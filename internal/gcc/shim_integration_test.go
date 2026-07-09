//go:build integration

package gcc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/EthanKim8683/cpx/internal/cdb"
	"github.com/EthanKim8683/cpx/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteGCC(t *testing.T) {
	t.Parallel()

	cfg, err := config.Load()
	if err != nil {
		t.Skipf("failed to load config: %v", err)
	}
	if cfg.GCC == "" {
		t.Skip("GCC not set")
	}

	dir, err := os.Getwd()
	require.NoError(t, err)
	tmpDir := t.TempDir()
	mainC := filepath.Join(tmpDir, "main.c")
	os.WriteFile(mainC, []byte("int main() { return 0; }"), 0644)
	mainO := filepath.Join(tmpDir, "main.o")
	cdbJSON := filepath.Join(tmpDir, ".cpx", "cdb.json")
	args := []string{
		gcc,
		"-o", mainO,
		"-std=c11",
		"-D", "DEFINE",
		"-O2",
		mainC,
	}
	command, err := cdb.Parse(CDBConfig, args)
	require.NoError(t, err)

	require.NoError(t, ExecuteGCC(&config.Config{
		GCC: cfg.GCC,
		CDB: cdbJSON,
	}, args))

	records, err := cdb.NewStore(cdbJSON).Records()
	require.NoError(t, err)
	require.ElementsMatch(t, []cdb.Record{
		{
			File:    mainC,
			Dir:     dir,
			Shim:    gcc,
			Command: command,
		},
	}, records)

	_, err = os.Stat(mainO)
	assert.NoError(t, err)
}

func TestExecuteGXX(t *testing.T) {
	t.Parallel()

	cfg, err := config.Load()
	if err != nil {
		t.Skipf("failed to load config: %v", err)
	}
	if cfg.GXX == "" {
		t.Skip("GXX not set")
	}

	dir, err := os.Getwd()
	require.NoError(t, err)
	tmpDir := t.TempDir()
	mainCPP := filepath.Join(tmpDir, "main.cpp")
	os.WriteFile(mainCPP, []byte("int main() { return 0; }"), 0644)
	mainO := filepath.Join(tmpDir, "main.o")
	cdbJSON := filepath.Join(tmpDir, ".cpx", "cdb.json")
	args := []string{
		gxx,
		"-o", mainO,
		"-std=c++17",
		"-D", "DEFINE",
		"-O2",
		mainCPP,
	}
	command, err := cdb.Parse(CDBConfig, args)
	require.NoError(t, err)

	require.NoError(t, ExecuteGXX(&config.Config{
		GXX: cfg.GXX,
		CDB: cdbJSON,
	}, args))

	records, err := cdb.NewStore(cdbJSON).Records()
	require.NoError(t, err)
	require.ElementsMatch(t, []cdb.Record{
		{
			File:    mainCPP,
			Dir:     dir,
			Shim:    gxx,
			Command: command,
		},
	}, records)

	_, err = os.Stat(mainO)
	assert.NoError(t, err)
}
