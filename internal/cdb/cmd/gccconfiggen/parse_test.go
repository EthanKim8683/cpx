package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNegateName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"fcommon", "fno-common"},
		{"Wextra", "Wno-extra"},
		{"msse", "mno-sse"},
		{"grecord", "gno-record"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, negateName(tt.input))
		})
	}
}

func TestNegateRE(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"fcommon", true},
		{"Wextra", true},
		{"msse", true},
		{"grecord", true},
		{"f", false},
		{"W", false},
		{"m", false},
		{"g", false},
		{"O3", false},
		{"fcommon=", false},
		{"Wno-all=always", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, negateRE.MatchString(tt.input))
		})
	}
}

func TestHasAttr(t *testing.T) {
	attrs := "Common RejectNegative Var(flag) Optimization"

	assert.True(t, hasAttr(attrs, "RejectNegative"))
	assert.True(t, hasAttr(attrs, "Common"))
	assert.False(t, hasAttr(attrs, "Var")) // Exact token match only
	assert.False(t, hasAttr(attrs, "Reject"))
}

func TestExpandOptRecords(t *testing.T) {
	records := []optRecord{
		{
			name:  "fcommon",
			attrs: "Common Var(flag_common)",
		},
		{
			name:  "Wextra",
			attrs: "Common RejectNegative",
		},
		{
			name:  "msse",
			attrs: "Target Alias(march, sse)",
		},
		{
			name:  "std=",
			attrs: "Joined",
		},
		{
			name:  "fextra",
			attrs: "Alias(fconserve-stack)",
		},
		{
			name:  "fother",
			attrs: "Alias(fconserve-stack, other_pos, other_neg)",
		},
	}

	got := expandOptRecords(records)

	expected := []optRecord{
		// Original records
		{
			name:  "fcommon",
			attrs: "Common Var(flag_common)",
		},
		{
			name:  "Wextra",
			attrs: "Common RejectNegative",
		},
		{
			name:  "msse",
			attrs: "Target Alias(march, sse)",
		},
		{
			name:  "std=",
			attrs: "Joined",
		},
		{
			name:  "fextra",
			attrs: "Alias(fconserve-stack)",
		},
		{
			name:  "fother",
			attrs: "Alias(fconserve-stack, other_pos, other_neg)",
		},
		// Synthesized negative records
		{
			name:  "fno-common",
			attrs: "RejectNegative",
		},
		{
			name:  "mno-sse",
			attrs: "RejectNegative",
		},
		{
			name:  "fno-extra",
			attrs: "RejectNegative",
		},
		{
			name:  "fno-other",
			attrs: "RejectNegative",
		},
	}

	assert.ElementsMatch(t, expected, got)
}
