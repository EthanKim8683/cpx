package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasProp(t *testing.T) {
	props := "Common RejectNegative Var(flag) Optimization"

	assert.True(t, hasProp("RejectNegative", props))
	assert.True(t, hasProp("Common", props))
	assert.True(t, hasProp("Var", props)) // Parentheses stripped, so "Var" is a separate word
	assert.False(t, hasProp("Reject", props))
}

func TestPropArgs(t *testing.T) {
	tests := []struct {
		name  string
		prop  string
		props string
		want  string
	}{
		{
			name:  "simple parameter",
			prop:  "Alias",
			props: "Common Alias(target)",
			want:  "target",
		},
		{
			name:  "multiple parameters",
			prop:  "Alias",
			props: "Common Alias(target, pos, neg)",
			want:  "target, pos, neg",
		},
		{
			name:  "braced parameters",
			prop:  "Alias",
			props: "Common Alias({target, pos})",
			want:  "target, pos",
		},
		{
			name:  "no match",
			prop:  "Alias",
			props: "Common RejectNegative",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, propArgs(tt.prop, tt.props))
		})
	}
}

func TestParseOptRecord(t *testing.T) {
	tests := []struct {
		name   string
		record optRecord
		want   parsedOptRecord
	}{
		{
			name: "flag option",
			record: optRecord{
				name:  "fcommon",
				props: "Common RejectNegative",
			},
			want: parsedOptRecord{
				name:           "fcommon",
				rejectNegative: true,
			},
		},
		{
			name: "joined option with Args(2)",
			record: optRecord{
				name:  "std=",
				props: "Joined Args(2)",
			},
			want: parsedOptRecord{
				name:   "std=",
				joined: true,
				args:   2,
			},
		},
		{
			name: "separate with NoDriverArg",
			record: optRecord{
				name:  "Q",
				props: "Separate NoDriverArg",
			},
			want: parsedOptRecord{
				name:        "Q",
				separate:    true,
				noDriverArg: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseOptRecord(tt.record)
			assert.Equal(t, tt.want, got)
		})
	}
}
