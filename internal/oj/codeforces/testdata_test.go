package codeforces

import (
	"bytes"
	"embed"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/*
var testdata embed.FS

func readDocument(t *testing.T, path string) *goquery.Document {
	t.Helper()

	html, err := testdata.ReadFile(path)
	require.NoError(t, err)

	d, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	require.NoError(t, err)
	return d
}
