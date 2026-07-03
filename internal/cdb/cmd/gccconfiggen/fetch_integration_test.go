//go:build integration

package main

import (
	"net/http"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

// TestFetchFileLiveIntegration performs a live network fetch from raw.githubusercontent.com
// for GCC's BASE-VER file, verifying that it contains the correct version string and matches the golden file.
func TestFetchFileLiveIntegration(t *testing.T) {
	version := "14.2.0"
	path := "gcc/BASE-VER"
	got, err := fetchFile(http.DefaultClient, version, path)
	require.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "BASE-VER", []byte(got))
}
