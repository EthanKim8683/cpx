package cbundler

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resolveTestdataRoot(t *testing.T) string {
	t.Helper()

	return filepath.Join(testdataRoot(t), "resolve")
}

func TestCommentIncludes(t *testing.T) {
	t.Parallel()

	source := `#include <foo.hpp>
#include "bar.hpp"
 # include"baz.hpp"
// #include "qux.hpp"`

	expected := `#include <foo.hpp>
// #include "bar.hpp"
//  # include"baz.hpp"
// #include "qux.hpp"`

	assert.Equal(t, expected, commentIncludes(source))
}

func TestFindIncludes(t *testing.T) {
	t.Parallel()

	source := `#include <foo.hpp>
#include "bar.hpp"
 # include"baz.hpp"
// #include "qux.hpp"`

	expected := []string{"bar.hpp", "baz.hpp"}
	assert.Equal(t, expected, findIncludes(source))
}

func TestResolveInclude(t *testing.T) {
	t.Parallel()

	r := resolveTestdataRoot(t)
	includePaths := []string{
		filepath.Join(r, "include"),
		filepath.Join(r, "include", "include"),
	}

	tests := map[string]struct {
		include string
		absPath string
		err     error
	}{
		"absolute include": {
			include: filepath.Join(r, "include", "foo.hpp"),
			absPath: filepath.Join(r, "include", "foo.hpp"),
			err:     nil,
		},
		"first include path": {
			include: "foo.hpp",
			absPath: filepath.Join(includePaths[0], "foo.hpp"),
			err:     nil,
		},
		"second include path": {
			include: "bar.hpp",
			absPath: filepath.Join(includePaths[1], "bar.hpp"),
			err:     nil,
		},
		"escaping include": {
			include: "../foo.hpp",
			absPath: filepath.Join(includePaths[0], "foo.hpp"),
			err:     nil,
		},
		"missing absolute include": {
			include: "/foo.hpp",
			absPath: "",
			err:     errors.New("could not resolve include: /foo.hpp"),
		},
		"missing include": {
			include: "qux.hpp",
			absPath: "",
			err:     errors.New("could not resolve include: qux.hpp"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			absPath, err := resolveInclude(test.include, includePaths)
			if test.err != nil {
				require.EqualError(t, err, test.err.Error())
				assert.Empty(t, absPath)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.absPath, absPath)
			}
		})
	}
}

func TestResolveIncludes(t *testing.T) {
	t.Parallel()

	r := resolveTestdataRoot(t)
	includePaths := []string{
		filepath.Join(r, "include"),
		filepath.Join(r, "include", "include"),
	}

	tests := map[string]struct {
		includes []string
		absPaths []string
		err      error
	}{
		"resolved all includes": {
			includes: []string{"foo.hpp", "bar.hpp"},
			absPaths: []string{
				filepath.Join(includePaths[0], "foo.hpp"),
				filepath.Join(includePaths[1], "bar.hpp"),
			},
		},
		"deduplicated includes": {
			includes: []string{"foo.hpp", "../foo.hpp"},
			absPaths: []string{
				filepath.Join(includePaths[0], "foo.hpp"),
			},
		},
		"missing some includes": {
			includes: []string{"foo.hpp", "/foo.hpp", "qux.hpp"},
			absPaths: []string{
				filepath.Join(includePaths[0], "foo.hpp"),
			},
			err: errors.Join(
				errors.New("could not resolve include: /foo.hpp"),
				errors.New("could not resolve include: qux.hpp"),
			),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			absPaths, err := resolveIncludes(test.includes, includePaths)
			if test.err != nil {
				require.EqualError(t, err, test.err.Error())
			} else {
				require.NoError(t, err)
			}
			assert.ElementsMatch(t, test.absPaths, absPaths)
		})
	}
}

func TestBuildIncludePaths(t *testing.T) {
	t.Parallel()

	absPath := filepath.Join("foo", "source.cpp")
	includePaths := []string{"bar", "baz"}

	collectedIncludePaths := buildIncludePaths(absPath, includePaths)
	assert.Equal(t, []string{"foo", "bar", "baz"}, collectedIncludePaths)
}

func TestResolveFile(t *testing.T) {
	t.Parallel()

	r := resolveTestdataRoot(t)
	includePaths := []string{
		filepath.Join(r, "include"),
		filepath.Join(r, "include", "include"),
	}
	var (
		sourcePath  = filepath.Join(r, "source.cpp")
		missingPath = filepath.Join(r, "missing.cpp")
	)

	t.Run("resolved file", func(t *testing.T) {
		t.Parallel()

		fragment, dependencies, err := resolveFile(sourcePath, includePaths)
		require.NoError(t, err)
		assert.Equal(t, `#include <foo.hpp>
// #include "bar.hpp"
//  # include"baz.hpp"
// #include "qux.hpp"`, fragment)
		assert.ElementsMatch(t, []string{
			filepath.Join(r, "bar.hpp"),
			filepath.Join(includePaths[0], "baz.hpp"),
		}, dependencies)
	})

	t.Run("missing file", func(t *testing.T) {
		t.Parallel()

		fragment, dependencies, err := resolveFile(missingPath, includePaths)
		require.ErrorIs(t, err, os.ErrNotExist)
		assert.Empty(t, fragment)
		assert.Empty(t, dependencies)
	})

	t.Run("missing some includes", func(t *testing.T) {
		t.Parallel()

		fragment, dependencies, err := resolveFile(sourcePath, []string{})
		require.EqualError(t, err, errors.New("could not resolve include: baz.hpp").Error())
		assert.Empty(t, fragment)
		assert.ElementsMatch(t, []string{
			filepath.Join(r, "bar.hpp"),
		}, dependencies)
	})
}
