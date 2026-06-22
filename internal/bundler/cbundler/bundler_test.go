package cbundler

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func bundlerTestdataRoot(t *testing.T) string {
	t.Helper()

	return filepath.Join(testdataRoot(t), "bundle")
}

func TestBundler(t *testing.T) {
	t.Parallel()

	r := bundlerTestdataRoot(t)

	tests := map[string]struct {
		sourcePath   string
		includePaths []string
		bundled      string
		err          error
	}{
		"valid": {
			sourcePath: filepath.Join(r, "valid", "foo.cpp"),
			includePaths: []string{
				filepath.Join(r, "valid", "include"),
			},
			bundled: `#include <iostream>

void bar() { std::cout << "Hello, World!" << std::endl; }
// #include "bar.hpp"

int main() { bar(); }`,
		},
		"relative": {
			sourcePath: filepath.Join("foo.cpp"),
			err:        errors.New("source path is not absolute: foo.cpp"),
		},
		"broken": {
			sourcePath: filepath.Join(r, "broken", "foo.cpp"),
			err:        errors.New("could not resolve include: bar.hpp"),
		},
		"cyclic": {
			sourcePath: filepath.Join(r, "cyclic", "foo.cpp"),
			err:        errors.New("cycle detected"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bundler := New(test.includePaths)
			bundled, err := bundler.Bundle(test.sourcePath)
			if test.err != nil {
				require.EqualError(t, err, test.err.Error())
				assert.Empty(t, bundled)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.bundled, bundled)
			}
		})
	}
}
