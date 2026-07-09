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

type mockRecordAdder struct {
	records []Record
	err     error
}

func (m *mockRecordAdder) Add(records []Record) error {
	m.records = append(m.records, records...)
	return m.err
}

func TestShim_Execute(t *testing.T) {
	t.Parallel()

	t.Run("successful execution", func(t *testing.T) {
		t.Parallel()
		compiler := &mockCompiler{}
		recordAdder := &mockRecordAdder{}

		shim := &Shim{
			Name:        "g++",
			Cfg:         &Config{Patterns: []OptionPattern{}},
			Compiler:    compiler,
			RecordAdder: recordAdder,
		}

		args := []string{"g++", "main.cpp"}
		err := shim.Execute(args)
		require.NoError(t, err)

		assert.Equal(t, args, compiler.argv)
		require.Len(t, recordAdder.records, 1)
		assert.Equal(t, "main.cpp", recordAdder.records[0].File)
	})

	t.Run("compiler failure", func(t *testing.T) {
		t.Parallel()
		compiler := &mockCompiler{err: errors.New("compilation failed")}
		recordAdder := &mockRecordAdder{}

		shim := &Shim{
			Name:        "g++",
			Cfg:         &Config{Patterns: []OptionPattern{}},
			Compiler:    compiler,
			RecordAdder: recordAdder,
		}

		args := []string{"g++", "main.cpp"}
		err := shim.Execute(args)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "compilation failed")

		// Record update still runs concurrently
		require.Len(t, recordAdder.records, 1)
	})

	t.Run("record adder failure", func(t *testing.T) {
		t.Parallel()
		compiler := &mockCompiler{}
		recordAdder := &mockRecordAdder{err: errors.New("db write failed")}

		shim := &Shim{
			Name:        "g++",
			Cfg:         &Config{Patterns: []OptionPattern{}},
			Compiler:    compiler,
			RecordAdder: recordAdder,
		}

		args := []string{"g++", "main.cpp"}
		err := shim.Execute(args)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "updating compilation database")
	})

	t.Run("empty arguments", func(t *testing.T) {
		t.Parallel()
		shim := &Shim{}
		err := shim.Execute([]string{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no compiler arguments provided")
	})
}
