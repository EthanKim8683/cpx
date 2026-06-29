package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestDetectVersion verifies that version detection executes against a system driver.
func TestDetectVersion(t *testing.T) {
	version, err := detectVersion("gcc")
	if err != nil {
		t.Logf("detectVersion('gcc') returned error (expected in environments without gcc): %v", err)
		return
	}
	if version == "" {
		t.Error("detectVersion('gcc') returned an empty version string")
	}
}

// TestRawURL verifies that GitHub raw download URLs are properly formatted for release tags.
func TestRawURL(t *testing.T) {
	got := rawURL("14.2.0", "gcc/common.opt")
	want := "https://raw.githubusercontent.com/gcc-mirror/gcc/releases/gcc-14.2.0/gcc/common.opt"
	if got != want {
		t.Errorf("rawURL() = %q; want %q", got, want)
	}
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
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("fetchSource() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("fetchSource() unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("content mismatch for %s (-want +got):\n%s", tt.path, diff)
			}
		})
	}
}

// TestFetchSourceLiveIntegration performs a live network fetch from raw.githubusercontent.com for a known GCC release file.
func TestFetchSourceLiveIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live network integration test in short mode")
	}

	version := "14.2.0"
	path := "gcc/common.opt"
	content, err := fetchSource(http.DefaultClient, version, path)
	if err != nil {
		t.Fatalf("fetchSource live integration failed for %s: %v", path, err)
	}
	if len(content) == 0 {
		t.Errorf("fetched in-memory file %s is empty", path)
	}
}

// newMockHTTPClient creates a test HTTP client that routes requests for raw.githubusercontent.com to a mock server.
func newMockHTTPClient(t *testing.T, handler http.HandlerFunc) *http.Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parsing mock server URL: %v", err)
	}

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
