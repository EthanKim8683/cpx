//go:build integration

package cdb

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreAddIntegration(t *testing.T) {
	tempDir := t.TempDir()
	dbFile := filepath.Join(tempDir, "cdb.json")
	store := NewStore(dbFile)

	records := []Record{
		{
			File: "main.cpp",
			Dir:  "/workspace",
			Shim: "g++",
		},
	}

	err := store.Add(records)
	require.NoError(t, err)

	// Verify the file exists and is populated
	_, err = os.Stat(dbFile)
	require.NoError(t, err)

	// Verify lock file exists
	_, err = os.Stat(dbFile + ".lock")
	require.NoError(t, err)

	data, err := os.ReadFile(dbFile)
	require.NoError(t, err)

	var stored []Record
	err = json.Unmarshal(data, &stored)
	require.NoError(t, err)
	require.Len(t, stored, 1)
	assert.Equal(t, "main.cpp", stored[0].File)
}
