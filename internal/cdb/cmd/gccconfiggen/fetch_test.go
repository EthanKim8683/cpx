package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

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

// TestFetchOptFileMock verifies file fetching over HTTP using a mock server.
func TestFetchOptFileMock(t *testing.T) {
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
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			path:    "gcc/common.opt",
			wantErr: errUnexpectedStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newMockHTTPClient(t, tt.handler)
			got, err := fetchFile(client, "14.2.0", tt.path)

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
