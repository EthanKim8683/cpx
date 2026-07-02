//go:build integration

package main

import (
	"net/http"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

// TestFetchOptFileLiveIntegration performs a live network fetch from raw.githubusercontent.com for a known GCC release file and verifies its content against a golden file.
func TestFetchOptFileLiveIntegration(t *testing.T) {
	version := "14.2.0"
	path := "gcc/analyzer/analyzer.opt"
	got, err := fetchOptFile(http.DefaultClient, version, path)
	require.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "analyzer.opt", []byte(got))
}
