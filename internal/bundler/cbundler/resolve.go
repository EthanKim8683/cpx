package cbundler

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var includeRegexp = regexp.MustCompile(`(?m)^\s*#\s*include\s*"((?:[^"]|\")+?)"`)

func commentIncludes(source string) string {
	return includeRegexp.ReplaceAllString(source, "// $0")
}

func findIncludes(source string) []string {
	matches := includeRegexp.FindAllStringSubmatch(source, -1)
	includes := make([]string, 0, len(matches))
	for _, match := range matches {
		includes = append(includes, match[1])
	}
	return includes
}

func resolveInclude(include string, includePaths []string) (string, error) {
	if filepath.IsAbs(include) {
		absPath := filepath.Clean(include)
		if _, err := os.Stat(absPath); err == nil {
			return absPath, nil
		}
	} else {
		for _, includePath := range includePaths {
			absPath := filepath.Clean(filepath.Join(includePath, include))
			if _, err := os.Stat(absPath); err == nil {
				return absPath, nil
			}
		}
	}
	return "", fmt.Errorf("could not resolve include: %s", include)
}

func resolveIncludes(includes []string, includePaths []string) ([]string, error) {
	absPathSet := make(map[string]struct{}, len(includes))
	var errs error
	for _, include := range includes {
		absPath, err := resolveInclude(include, includePaths)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		absPathSet[absPath] = struct{}{}
	}

	absPaths := make([]string, 0, len(absPathSet))
	for absPath := range absPathSet {
		absPaths = append(absPaths, absPath)
	}
	return absPaths, errs
}

func buildIncludePaths(absPath string, includePaths []string) []string {
	return append([]string{filepath.Dir(absPath)}, includePaths...)
}

func resolveFile(absPath string, includePaths []string) (string, []string, error) {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", nil, err
	}
	source := string(data)

	dependencies, err := resolveIncludes(
		findIncludes(source),
		buildIncludePaths(absPath, includePaths),
	)
	if err != nil {
		return "", dependencies, err
	}

	return commentIncludes(source), dependencies, nil
}
