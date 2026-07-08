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

func TestShim_Integration(t *testing.T) {
	t.Parallel()

	t.Run("update writes to store", func(t *testing.T) {
		t.Parallel()

		tempDir := t.TempDir()
		dbFile := filepath.Join(tempDir, "cdb.json")
		store := NewStore(dbFile)

		shim := &Shim{
			Name:  "g++",
			Bin:   "echo", // not used by update()
			Cfg:   &Config{Patterns: []OptionPattern{}},
			Store: store,
		}

		args := []string{"g++", "main.cpp", "solve.cpp"}
		err := shim.update(args)
		require.NoError(t, err)

		// Verify db file exists and has correct records
		data, err := os.ReadFile(dbFile)
		require.NoError(t, err)

		var stored []Record
		err = json.Unmarshal(data, &stored)
		require.NoError(t, err)
		assert.Len(t, stored, 2)
		
		files := map[string]bool{
			stored[0].File: true,
			stored[1].File: true,
		}
		assert.True(t, files["main.cpp"])
		assert.True(t, files["solve.cpp"])
	})

	t.Run("compile runs binary", func(t *testing.T) {
		t.Parallel()

		shim := &Shim{
			Name: "g++",
			Bin:  "true", // exits with 0 immediately
		}

		args := []string{"g++"}
		err := shim.compile(args)
		require.NoError(t, err)
	})

	t.Run("Execute runs both compile and update", func(t *testing.T) {
		t.Parallel()

		tempDir := t.TempDir()
		dbFile := filepath.Join(tempDir, "cdb.json")
		store := NewStore(dbFile)

		shim := &Shim{
			Name:  "g++",
			Bin:   "true", // exits with 0 immediately
			Cfg:   &Config{Patterns: []OptionPattern{}},
			Store: store,
		}

		args := []string{"g++", "main.cpp"}
		err := shim.Execute(args)
		require.NoError(t, err)

		// Verify compiler output database was created
		_, err = os.Stat(dbFile)
		require.NoError(t, err)
	})
}
