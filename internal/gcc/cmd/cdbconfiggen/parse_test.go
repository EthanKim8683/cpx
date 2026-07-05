package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsOptRecord(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "option record",
			content: "O\nJoined",
			want:    true,
		},
		{
			name:    "excluded Language record",
			content: "Language\nC",
			want:    false,
		},
		{
			name:    "excluded Variable record",
			content: "Variable\nint target_flags",
			want:    false,
		},
		{
			name:    "excluded Enum record",
			content: "Enum\nValue",
			want:    false,
		},
		{
			name:    "empty content",
			content: "",
			want:    false,
		},
		{
			name:    "whitespace only",
			content: "  \n  ",
			want:    false,
		},
		{
			name:    "excluded TargetVariable",
			content: "TargetVariable\nint foo",
			want:    false,
		},
		{
			name:    "excluded TargetSave",
			content: "TargetSave\nint foo",
			want:    false,
		},
		{
			name:    "excluded HeaderKeep",
			content: "HeaderKeep\nfoo.h",
			want:    false,
		},
		{
			name:    "excluded EnumValue",
			content: "EnumValue\nFoo BAR 0",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isOptRecord(tt.content)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseOptRecord(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantName string
		wantProps string
	}{
		{
			name:      "flag option",
			content:   "static\nDriver",
			wantName:  "static",
			wantProps: "Driver",
		},
		{
			name:      "joined option",
			content:   "std=\nJoined Common",
			wantName:  "std=",
			wantProps: "Joined Common",
		},
		{
			name:      "no properties",
			content:   "foo",
			wantName:  "foo",
			wantProps: "",
		},
		{
			name:      "extra lines beyond second are ignored",
			content:   "foo\nJoined\nThis is description text",
			wantName:  "foo",
			wantProps: "Joined",
		},
		{
			name:      "leading whitespace trimmed",
			content:   "  foo  \n  Joined  ",
			wantName:  "foo",
			wantProps: "Joined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseOptRecord(tt.content)
			assert.Equal(t, tt.wantName, got.name)
			assert.Equal(t, tt.wantProps, got.props)
		})
	}
}

func TestParseOptRecords(t *testing.T) {
	t.Run("single option record", func(t *testing.T) {
		content := "foo\nJoined Driver"
		got := parseOptRecords(content)
		assert.Len(t, got, 1)
		assert.Equal(t, "foo", got[0].name)
		assert.Equal(t, "Joined Driver", got[0].props)
	})

	t.Run("multiple records separated by blank lines", func(t *testing.T) {
		content := "foo\nJoined\n\nbar\nSeparate"
		got := parseOptRecords(content)
		assert.Len(t, got, 2)
	})

	t.Run("comments are stripped", func(t *testing.T) {
		content := "; this is a comment\nfoo\nJoined\n; another comment"
		got := parseOptRecords(content)
		assert.Len(t, got, 1)
		assert.Equal(t, "foo", got[0].name)
	})

	t.Run("excluded record types are skipped", func(t *testing.T) {
		content := "Variable\nint target_flags\n\nfoo\nJoined"
		got := parseOptRecords(content)
		assert.Len(t, got, 1)
		assert.Equal(t, "foo", got[0].name)
	})

	t.Run("CRLF line endings are normalized", func(t *testing.T) {
		content := "foo\r\nJoined\r\n\r\nbar\r\nSeparate"
		got := parseOptRecords(content)
		assert.Len(t, got, 2)
	})

	t.Run("empty input", func(t *testing.T) {
		got := parseOptRecords("")
		assert.Empty(t, got)
	})

	t.Run("only comments and blank lines", func(t *testing.T) {
		content := "; comment\n\n; another comment\n"
		got := parseOptRecords(content)
		assert.Empty(t, got)
	})

	t.Run("mixed records and exclusions", func(t *testing.T) {
		content := "; Header\nVariable\nint flags\n\nfoo\nJoined RejectDriver\n\nbar\nSeparate"
		got := parseOptRecords(content)
		assert.Len(t, got, 2)
		assert.Equal(t, "foo", got[0].name)
		assert.Equal(t, "bar", got[1].name)
	})
}
