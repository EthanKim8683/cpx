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
)

// errUnexpectedStatus is returned when a remote fetch yields a non-200 HTTP response.
var errUnexpectedStatus = errors.New("unexpected HTTP status")

// fetchFile downloads a single GCC option specification file from the git mirror.
// It returns the file content as a string, or an error if the request fails or returns a non-200 status.
func fetchFile(client *http.Client, baseURL, path string) (string, error) {
	u := fmt.Sprintf("%s/%s", baseURL, path)
	resp, err := client.Get(u)
	if err != nil {
		return "", fmt.Errorf("fetching %s: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetching %s: %w", u, errUnexpectedStatus)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading body of %s: %w", u, err)
	}
	return string(b), nil
}
