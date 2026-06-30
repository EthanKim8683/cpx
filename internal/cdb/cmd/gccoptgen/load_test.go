package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDetectVersion verifies that version detection executes against a system driver.
func TestDetectVersion(t *testing.T) {
	version, err := detectVersion("gcc")
	if err != nil {
		t.Logf("detectVersion('gcc') returned error (expected in environments without gcc): %v", err)
		return
	}
	assert.NotEmpty(t, version, "detectVersion('gcc') returned an empty version string")
}

// TestRawURL verifies that GitHub raw download URLs are properly formatted for release tags.
func TestRawURL(t *testing.T) {
	got := rawURL("14.2.0", "gcc/common.opt")
	want := "https://raw.githubusercontent.com/gcc-mirror/gcc/releases/gcc-14.2.0/gcc/common.opt"
	assert.Equal(t, want, got)
}

// TestFetchSourceMock verifies file fetching over HTTP using a mock server.
func TestFetchSourceMock(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		path    string
		want    string
		wantErr error
	}{
		{
			name: "successful fetch",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "mock content for %s", r.URL.Path)
			},
			path: "gcc/common.opt",
			want: "mock content for /gcc-mirror/gcc/releases/gcc-14.2.0/gcc/common.opt",
		},
		{
			name: "unexpected 404 status",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			path:    "gcc/common.opt",
			wantErr: ErrUnexpectedStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newMockHTTPClient(t, tt.handler)
			got, err := fetchSource(client, "14.2.0", tt.path)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// newMockHTTPClient creates a test HTTP client that routes requests for raw.githubusercontent.com to a mock server.
func newMockHTTPClient(t *testing.T, handler http.HandlerFunc) *http.Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	serverURL, err := url.Parse(server.URL)
	require.NoError(t, err)

	client := server.Client()
	origTransport := client.Transport
	client.Transport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		req.URL.Scheme = serverURL.Scheme
		req.URL.Host = serverURL.Host
		return origTransport.RoundTrip(req)
	})

	return client
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// TestLoadMock verifies that load successfully runs detectVersion, downloads all required GCC option files, and creates the directory structure.
func TestLoadMock(t *testing.T) {
	// 1. Create a dummy compiler script that outputs a dummy GCC version on the real OS filesystem (needed for exec.Command).
	compilerScriptContent := "#!/bin/sh\necho 14.2.0\n"
	tmpCompiler, err := os.CreateTemp("", "mock-gcc-*")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Remove(tmpCompiler.Name()) })

	_, err = tmpCompiler.WriteString(compilerScriptContent)
	require.NoError(t, err)
	err = tmpCompiler.Close()
	require.NoError(t, err)
	err = os.Chmod(tmpCompiler.Name(), 0755)
	require.NoError(t, err)

	// 2. Setup mock HTTP server for the fetches.
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "mock content for %s", r.URL.Path)
	}
	client := newMockHTTPClient(t, handler)

	// 3. Run load using an in-memory Afero Fs.
	mockFS := afero.NewMemMapFs()
	afs := &afero.Afero{Fs: mockFS}
	tempDir := "/temp"
	version, err := load(mockFS, client, tmpCompiler.Name(), tempDir)
	require.NoError(t, err)
	assert.Equal(t, "14.2.0", version)

	// 4. Assert the exact file structure and contents match using Walk + Testify map assertion.
	gotFiles := make(map[string]string)
	err = afero.Walk(mockFS, tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(tempDir, path)
		if err != nil {
			return err
		}
		content, err := afs.ReadFile(path)
		if err != nil {
			return err
		}
		gotFiles[filepath.ToSlash(rel)] = string(content)
		return nil
	})
	require.NoError(t, err)

	wantFiles := make(map[string]string)
	for _, file := range files {
		wantFiles[file] = "mock content for /gcc-mirror/gcc/releases/gcc-14.2.0/" + file
	}

	assert.Equal(t, wantFiles, gotFiles)
}
