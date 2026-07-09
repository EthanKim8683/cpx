package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

// findOptFiles recursively searches the given directory in the filesystem fs
// for files with a ".opt" extension. It collects and joins any errors encountered
// during the traversal.
func findOptFiles(fs afero.Fs, dir string) ([]string, error) {
	var files []string
	var errs error
	//nolint:errcheck // Walk errors are handled and joined via the inner callback
	_ = afero.Walk(fs, dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to walk %s: %w", path, err))
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".opt" {
			files = append(files, path)
		}
		return nil
	})
	if errs != nil {
		return nil, errs
	}
	return files, nil
}

// readOptFiles reads the content of each file in files from the filesystem fs.
// It returns a slice containing the contents of each file in the same order.
// If any file read fails, it continues reading the remaining files and returns
// all encountered errors joined together.
func readOptFiles(fs afero.Fs, files []string) ([]string, error) {
	contents := make([]string, 0, len(files))
	var errs error
	for _, file := range files {
		content, err := afero.ReadFile(fs, file)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to read file %s: %w", file, err))
			continue
		}
		contents = append(contents, string(content))
	}
	if errs != nil {
		return nil, errs
	}
	return contents, nil
}
