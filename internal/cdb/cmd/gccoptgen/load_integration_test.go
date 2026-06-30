//go:build integration

package main

import (
	"net/http"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

// TestFetchSourceLiveIntegration performs a live network fetch from raw.githubusercontent.com for a known GCC release file and verifies its content against a golden file.
func TestFetchSourceLiveIntegration(t *testing.T) {
	version := "14.2.0"
	path := "gcc/opt-gather.awk"
	got, err := fetchSource(http.DefaultClient, version, path)
	require.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "opt-gather.awk", []byte(got))
}
