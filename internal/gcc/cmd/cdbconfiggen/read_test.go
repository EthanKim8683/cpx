package main

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindOptFiles(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		files := []string{
			"a.ext",
			"b.opt",
			"c/d.opt",
			"e/f.ext",
			"c/g/h/i/j.opt",
		}
		want := []string{
			"c/d.opt",
			"c/g/h/i/j.opt",
		}

		fs := afero.NewMemMapFs()
		for _, file := range files {
			err := fs.MkdirAll(filepath.Dir(file), 0o755)
			require.NoError(t, err)
			f, err := fs.Create(file)
			require.NoError(t, err)
			err = f.Close()
			require.NoError(t, err)
		}

		got, err := findOptFiles(fs, "c")
		require.NoError(t, err)
		assert.ElementsMatch(t, want, got)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		// Walk will fail if the directory does not exist in MemMapFs
		_, err := findOptFiles(fs, "non_existent_dir")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to walk non_existent_dir")
	})
}

func TestReadOptFiles(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		files := []string{"file1.opt", "file2.opt"}
		contents := []string{"content1", "content2"}

		for i, file := range files {
			f, err := fs.Create(file)
			require.NoError(t, err)
			_, err = f.Write([]byte(contents[i]))
			require.NoError(t, err)
			err = f.Close()
			require.NoError(t, err)
		}

		got, err := readOptFiles(fs, files)
		require.NoError(t, err)
		assert.Equal(t, contents, got)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		// Attempt to read files that do not exist to trigger and test error joining
		_, err := readOptFiles(fs, []string{"non_existent1.opt", "non_existent2.opt"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read file non_existent1.opt")
		assert.Contains(t, err.Error(), "failed to read file non_existent2.opt")
	})
}
