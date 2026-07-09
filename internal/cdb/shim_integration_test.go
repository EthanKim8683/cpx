//go:build integration

package cdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRecordAdder struct {
	records []Record
	err     error
}

func (m *mockRecordAdder) Add(records []Record) error {
	m.records = append(m.records, records...)
	return m.err
}

func TestShim_Integration(t *testing.T) {
	t.Parallel()

	t.Run("update writes to store", func(t *testing.T) {
		t.Parallel()

		store := &mockRecordAdder{}
		shim := &Shim{
			Name:        "g++",
			Cfg:         &Config{Patterns: []OptionPattern{}},
			Compiler:    &ExecCompiler{Bin: "echo"},
			RecordAdder: store,
		}

		args := []string{"g++", "main.cpp", "solve.cpp"}
		err := shim.update(args)
		require.NoError(t, err)

		require.Len(t, store.records, 2)
		assert.Equal(t, "main.cpp", store.records[0].File)
		assert.Equal(t, "solve.cpp", store.records[1].File)
	})

	t.Run("Execute runs both compile and update", func(t *testing.T) {
		t.Parallel()

		store := &mockRecordAdder{}
		shim := &Shim{
			Name:        "g++",
			Cfg:         &Config{Patterns: []OptionPattern{}},
			Compiler:    &ExecCompiler{Bin: "true"},
			RecordAdder: store,
		}

		args := []string{"g++", "main.cpp"}
		err := shim.Execute(args)
		require.NoError(t, err)

		require.Len(t, store.records, 1)
		assert.Equal(t, "main.cpp", store.records[0].File)
	})
}
