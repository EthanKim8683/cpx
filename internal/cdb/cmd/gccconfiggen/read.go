package main

import (
	"errors"
	"fmt"

	"github.com/spf13/afero"
)

// readOptFiles reads all .opt files in the given directory and returns their contents.
func readOptFiles(dir string) ([]string, error) {
	fs := afero.NewBasePathFs(afero.NewOsFs(), dir)
	files, err := afero.Glob(fs, "**/*.opt")
	if err != nil {
		return nil, fmt.Errorf("failed to glob files in %s: %w", dir, err)
	}

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
