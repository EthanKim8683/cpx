package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsOptRecord(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "valid option",
			content:  "fcommon\nCommon Var(flag_conserve_stack)",
			expected: true,
		},
		{
			name:     "excluded option",
			content:  "Variable\nint flag_var = 0",
			expected: false,
		},
		{
			name:     "empty content",
			content:  "",
			expected: false,
		},
		{
			name:     "only whitespace content",
			content:  "   \n  ",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isOptRecord(tt.content))
		})
	}
}

func TestExtractOptRecord(t *testing.T) {
	t.Run("with attributes", func(t *testing.T) {
		content := "fcommon\nCommon Var(flag_conserve_stack)\nSome description here"
		got := extractOptRecord(content)
		assert.Equal(t, "fcommon", got.name)
		assert.Equal(t, "Common Var(flag_conserve_stack)", got.attrs)
	})

	t.Run("without attributes", func(t *testing.T) {
		content := "fcommon"
		got := extractOptRecord(content)
		assert.Equal(t, "fcommon", got.name)
		assert.Empty(t, got.attrs)
	})
}

func TestExtractOptRecords(t *testing.T) {
	content := `; This is a comment
; Another comment

fcommon
Common Var(flag_conserve_stack) Init(0)
Enable normal stack conservation.

Variable
int flag_var = 0

fother
Common Joined
Another option.
`

	records := extractOptRecords(content)

	expected := []optRecord{
		{
			name:  "fcommon",
			attrs: "Common Var(flag_conserve_stack) Init(0)",
		},
		{
			name:  "fother",
			attrs: "Common Joined",
		},
	}

	assert.Equal(t, expected, records)
}
