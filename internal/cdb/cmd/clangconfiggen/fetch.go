package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/spf13/afero"
)

// errUnexpectedStatus is returned when a remote fetch yields a non-200 HTTP response.
var errUnexpectedStatus = errors.New("unexpected HTTP status")

// fetchFile fetches a source file from a remote URL and writes it to a local filesystem.
func fetchFile(client *http.Client, fs afero.Fs, baseURL, path string) error {
	u := fmt.Sprintf("%s/%s", baseURL, path)
	resp, err := client.Get(u)
	if err != nil {
		return fmt.Errorf("fetching %s: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetching %s: %w", u, errUnexpectedStatus)
	}

	afs := afero.Afero{Fs: fs}
	if err := afs.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating parent dirs for %s: %w", path, err)
	}

	f, err := afs.Create(path)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", path, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("writing %s: %w", u, err)
	}
	return nil
}
