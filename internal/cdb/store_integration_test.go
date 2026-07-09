//go:build integration

package cdb

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileRecordAdder_Add(t *testing.T) {
	t.Parallel()

	t.Run("successfully add record", func(t *testing.T) {
		t.Parallel()

		tempDir := t.TempDir()
		dbFile := filepath.Join(tempDir, "cdb.json")
		store := NewFileRecordAdder(dbFile)

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
	})

	t.Run("handling corrupt JSON", func(t *testing.T) {
		t.Parallel()

		tempDir := t.TempDir()
		dbFile := filepath.Join(tempDir, "cdb.json")
		store := NewFileRecordAdder(dbFile)

		// Write corrupt JSON to the database file
		err := os.WriteFile(dbFile, []byte("{not valid json"), 0644)
		require.NoError(t, err)

		records := []Record{
			{File: "main.cpp", Dir: "/workspace", Shim: "g++"},
		}

		err = store.Add(records)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "parsing database JSON")
	})

	t.Run("handling empty records", func(t *testing.T) {
		t.Parallel()

		tempDir := t.TempDir()
		dbFile := filepath.Join(tempDir, "cdb.json")
		store := NewFileRecordAdder(dbFile)

		// Add empty records — should succeed without error
		err := store.Add([]Record{})
		require.NoError(t, err)

		// Database file should exist but contain empty JSON array
		data, err := os.ReadFile(dbFile)
		require.NoError(t, err)

		var stored []Record
		require.NoError(t, json.Unmarshal(data, &stored))
		assert.Empty(t, stored)
	})

	t.Run("concurrent updates", func(t *testing.T) {
		t.Parallel()

		tempDir := t.TempDir()
		dbFile := filepath.Join(tempDir, "cdb.json")
		store := NewFileRecordAdder(dbFile)

		const goroutines = 50
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := range goroutines {
			go func() {
				defer wg.Done()
				records := []Record{
					{
						File: fmt.Sprintf("file_%d.cpp", i),
						Dir:  fmt.Sprintf("/dir_%d", i),
						Shim: "g++",
					},
				}
				err := store.Add(records)
				assert.NoError(t, err)
			}()
		}

		wg.Wait()

		data, err := os.ReadFile(dbFile)
		require.NoError(t, err)

		var stored []Record
		require.NoError(t, json.Unmarshal(data, &stored))
		require.Len(t, stored, goroutines)

		byFile := make(map[string]Record, len(stored))
		for _, r := range stored {
			byFile[r.File] = r
		}

		for i := range goroutines {
			key := fmt.Sprintf("file_%d.cpp", i)
			r, ok := byFile[key]
			require.True(t, ok, "missing record for %s", key)
			assert.Equal(t, fmt.Sprintf("/dir_%d", i), r.Dir)
		}
	})

	t.Run("concurrent overwrite matches key", func(t *testing.T) {
		t.Parallel()

		tempDir := t.TempDir()
		dbFile := filepath.Join(tempDir, "cdb.json")
		store := NewFileRecordAdder(dbFile)

		const goroutines = 50
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := range goroutines {
			go func() {
				defer wg.Done()
				records := []Record{
					{
						File: "shared.cpp",
						Dir:  fmt.Sprintf("/dir_%d", i),
						Shim: "g++",
					},
				}
				err := store.Add(records)
				assert.NoError(t, err)
			}()
		}

		wg.Wait()

		data, err := os.ReadFile(dbFile)
		require.NoError(t, err)

		var stored []Record
		require.NoError(t, json.Unmarshal(data, &stored))
		require.Len(t, stored, 1, "shared.cpp should exist exactly once")
		assert.Equal(t, "shared.cpp", stored[0].File)
	})
}
