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

	tmpDir := t.TempDir()
	mainC := filepath.Join(tmpDir, "main.c")
	mainO := filepath.Join(tmpDir, "main.o")
	cdbJSON := filepath.Join(tmpDir, ".cpx", "cdb.json")

	os.WriteFile(mainC, []byte("int main() { return 0; }"), 0644)

	args := []string{
		gcc,
		"-o", mainO,
		"-std=c11",
		"-D", "DEFINE",
		"-O2",
		mainC,
	}
	require.NoError(t, ExecuteGCC(&config.Config{
		GCC: cfg.GCC,
		CDB: cdbJSON,
	}, args, tmpDir))

	records, err := cdb.NewFileStore(cdbJSON).Records()
	require.NoError(t, err)

	command, err := cdb.Parse(CDBConfig, args)
	require.NoError(t, err)
	assert.ElementsMatch(t, []cdb.Record{
		{
			File:    mainC,
			Dir:     tmpDir,
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

	tmpDir := t.TempDir()
	mainCPP := filepath.Join(tmpDir, "main.cpp")
	mainO := filepath.Join(tmpDir, "main.o")
	cdbJSON := filepath.Join(tmpDir, ".cpx", "cdb.json")

	os.WriteFile(mainCPP, []byte("int main() { return 0; }"), 0644)

	args := []string{
		gxx,
		"-o", mainO,
		"-std=c11",
		"-D", "DEFINE",
		"-O2",
		mainCPP,
	}
	require.NoError(t, ExecuteGXX(&config.Config{
		GXX: cfg.GXX,
		CDB: cdbJSON,
	}, args, tmpDir))

	records, err := cdb.NewFileStore(cdbJSON).Records()
	require.NoError(t, err)

	command, err := cdb.Parse(CDBConfig, args)
	require.NoError(t, err)
	assert.ElementsMatch(t, []cdb.Record{
		{
			File:    mainCPP,
			Dir:     tmpDir,
			Shim:    gxx,
			Command: command,
		},
	}, records)

	_, err = os.Stat(mainO)
	assert.NoError(t, err)
}
