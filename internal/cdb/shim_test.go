package cdb

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCompiler struct {
	argv []string
	err  error
}

func (m *mockCompiler) Compile(argv []string) error {
	m.argv = argv
	return m.err
}

type mockStore struct {
	records []Record
	err     error
}

func (m *mockStore) Add(records []Record) error {
	m.records = append(m.records, records...)
	return m.err
}

func (m *mockStore) Records() ([]Record, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.records, nil
}

func TestShim_Execute(t *testing.T) {
	t.Parallel()

	t.Run("successful execution", func(t *testing.T) {
		t.Parallel()
		compiler := &mockCompiler{}
		store := &mockStore{}

		name := "shim"
		cfg := &Config{Patterns: []OptionPattern{}}
		shim := &Shim{
			Name:     name,
			Cfg:      cfg,
			Compiler: compiler,
			Store:    store,
		}

		file := "file"
		args := []string{"compiler", file}
		dir := "dir"
		err := shim.Execute(args, dir)
		require.NoError(t, err)

		records, err := store.Records()
		require.NoError(t, err)

		command, err := Parse(cfg, args)
		require.NoError(t, err)
		assert.ElementsMatch(t, []Record{
			{
				File:    file,
				Dir:     dir,
				Shim:    name,
				Command: command,
			},
		}, records)
	})

	t.Run("compiler failure", func(t *testing.T) {
		t.Parallel()
		compiler := &mockCompiler{err: errors.New("compilation failed")}
		store := &mockStore{}

		name := "shim"
		cfg := &Config{Patterns: []OptionPattern{}}
		shim := &Shim{
			Name:     name,
			Cfg:      cfg,
			Compiler: compiler,
			Store:    store,
		}

		file := "file"
		args := []string{"compiler", file}
		dir := "dir"
		err := shim.Execute(args, dir)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "compilation failed")

		// Store update still runs concurrently with compilation.
		records, err := store.Records()
		require.NoError(t, err)

		command, err := Parse(cfg, args)
		require.NoError(t, err)
		assert.ElementsMatch(t, []Record{
			{
				File:    file,
				Dir:     dir,
				Shim:    name,
				Command: command,
			},
		}, records)
	})

	t.Run("store failure", func(t *testing.T) {
		t.Parallel()
		compiler := &mockCompiler{}
		store := &mockStore{err: errors.New("db write failed")}

		shim := &Shim{
			Name:     "shim",
			Cfg:      &Config{Patterns: []OptionPattern{}},
			Compiler: compiler,
			Store:    store,
		}

		err := shim.Execute([]string{"compiler"}, "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "updating compilation database")
	})

	t.Run("empty arguments", func(t *testing.T) {
		t.Parallel()
		shim := &Shim{}
		err := shim.Execute([]string{}, "")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no compiler arguments provided")
	})
}
