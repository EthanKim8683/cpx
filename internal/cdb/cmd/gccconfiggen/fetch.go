// fetch.go provides functions to query the local compiler executable version
// and retrieve remote option specification (.opt) files from the GCC mirror repository.
// For details on GCC options, see the GCC Command Options Documentation
// (https://gcc.gnu.org/onlinedocs/gcc/Option-Summary.html).

package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

// errUnexpectedStatus is returned when a remote fetch yields a non-200 HTTP response.
var errUnexpectedStatus = errors.New("unexpected HTTP status")

// detectVersion queries the GCC driver binary at path to extract its release version string.
// The path must be non-empty and point to a valid GCC binary.
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

// rawURL resolves the GitHub raw download URL for a GCC source file at a given version tag.
func rawURL(version, path string) string {
	tag := fmt.Sprintf("releases/gcc-%s", version)
	return fmt.Sprintf("https://raw.githubusercontent.com/gcc-mirror/gcc/%s/%s", tag, path)
}

// fetchFile downloads a single GCC option specification file from the git mirror.
// It returns the file content as a string, or an error if the request fails or returns a non-200 status.
func fetchFile(client *http.Client, version, path string) (string, error) {
	u := rawURL(version, path)
	resp, err := client.Get(u)
	if err != nil {
		return "", fmt.Errorf("fetching file %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetching file %s from %s: %w (%d)", path, u, errUnexpectedStatus, resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading file body of %s: %w", path, err)
	}
	return string(b), nil
}
