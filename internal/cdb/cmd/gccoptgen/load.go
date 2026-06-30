// Package main implements the GCC option config generator tool.
package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// ErrUnexpectedStatus indicates an HTTP request returned a non-200 status code.
var ErrUnexpectedStatus = errors.New("unexpected HTTP status")

// detectVersion queries the GCC driver binary at path to extract its release version string (e.g., "14.2.0").
func detectVersion(path string) (string, error) {
	// GCC 7 introduced -dumpfullversion to guarantee a 3-part version string (major.minor.patch) suitable
	// for release tag matching (https://gcc.gnu.org/gcc-7/changes.html).
	cmd := exec.Command(path, "-dumpfullversion")
	out, err := cmd.Output()
	if err != nil {
		// Compilers older than GCC 7 do not support -dumpfullversion, but -dumpversion returned the full version on those releases.
		cmd = exec.Command(path, "-dumpversion")
		out, err = cmd.Output()
		if err != nil {
			return "", fmt.Errorf("detecting GCC version via %s: %w", path, err)
		}
	}

	version := strings.TrimSpace(string(out))
	if version == "" {
		return "", fmt.Errorf("detecting GCC version via %s: empty output returned", path)
	}

	return version, nil
}

// rawURL constructs the raw HTTP download URL for a GCC repository file at a specific release version tag.
func rawURL(version, path string) string {
	tag := fmt.Sprintf("releases/gcc-%s", version)
	return fmt.Sprintf("https://raw.githubusercontent.com/gcc-mirror/gcc/%s/%s", tag, path)
}

// fetchSource downloads a single GCC repository file over HTTP for a specified release version.
// It returns the raw text content string of the requested file relative path (e.g. "gcc/common.opt").
func fetchSource(client *http.Client, version, path string) (string, error) {
	u := rawURL(version, path)
	resp, err := client.Get(u)
	if err != nil {
		return "", fmt.Errorf("fetching source %s: %w", path, err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return "", fmt.Errorf("fetching source %s from %s: %w (%d)", path, u, ErrUnexpectedStatus, resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("reading source body of %s: %w", path, err)
	}

	return string(b), nil
}

// files is the manifest of upstream GCC files required for option generation.
var files = []string{
	"gcc/common.opt",
	"gcc/c-family/c.opt",
	"gcc/params.opt",
	"gcc/analyzer/analyzer.opt",
	"gcc/opt-functions.awk",
	"gcc/opt-read.awk",
	"gcc/opt-gather.awk",
}

// load detects the compiler version, fetches all required source option files, and writes them to the specified directory.
// It returns the detected compiler version, or an error if the process fails.
func load(fs afero.Fs, client *http.Client, compilerPath, dir string) (string, error) {
	version, err := detectVersion(compilerPath)
	if err != nil {
		return "", fmt.Errorf("detecting version: %w", err)
	}

	afs := &afero.Afero{Fs: fs}
	for _, file := range files {
		content, err := fetchSource(client, version, file)
		if err != nil {
			return "", fmt.Errorf("fetching source %s: %w", file, err)
		}

		destPath := filepath.Join(dir, filepath.FromSlash(file))
		if err := afs.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return "", fmt.Errorf("creating parent directories for %s: %w", destPath, err)
		}

		if err := afs.WriteFile(destPath, []byte(content), 0644); err != nil {
			return "", fmt.Errorf("writing file %s: %w", destPath, err)
		}
	}

	return version, nil
}

