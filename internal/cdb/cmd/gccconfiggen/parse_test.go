package main

import (
	"testing"

	"github.com/EthanKim8683/cpx/internal/cdb"
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

func TestParseOptRecord(t *testing.T) {
	tests := []struct {
		name       string
		record     optRecord
		wantSpell  string
		wantKind   cdb.OptionKind
		wantNumArg int
	}{
		{
			name: "flag option",
			record: optRecord{
				name:  "fcommon",
				attrs: "Common Var(flag_common)",
			},
			wantSpell:  "-fcommon",
			wantKind:   cdb.OptionKindFlag,
			wantNumArg: 0,
		},
		{
			name: "joined option",
			record: optRecord{
				name:  "std=",
				attrs: "Joined RejectNegative",
			},
			wantSpell:  "-std=",
			wantKind:   cdb.OptionKindJoined,
			wantNumArg: 1,
		},
		{
			name: "separate option",
			record: optRecord{
				name:  "o",
				attrs: "Separate",
			},
			wantSpell:  "-o",
			wantKind:   cdb.OptionKindSeparate,
			wantNumArg: 1,
		},
		{
			name: "joined or separate option",
			record: optRecord{
				name:  "I",
				attrs: "Joined Separate",
			},
			wantSpell:  "-I",
			wantKind:   cdb.OptionKindJoinedOrSeparate,
			wantNumArg: 1,
		},
		{
			name: "comma joined option",
			record: optRecord{
				name:  "fsanitize=",
				attrs: "CommaJoined Joined",
			},
			wantSpell:  "-fsanitize=",
			wantKind:   cdb.OptionKindCommaJoined,
			wantNumArg: 1,
		},
		{
			name: "multi-arg option (Args(2))",
			record: optRecord{
				name:  "sectcreate",
				attrs: "Args(2) Separate",
			},
			wantSpell:  "-sectcreate",
			wantKind:   cdb.OptionKindMultiArg,
			wantNumArg: 2,
		},
		{
			name: "multi-arg option (Args(4))",
			record: optRecord{
				name:  "fourargs",
				attrs: "Args(4)",
			},
			wantSpell:  "-fourargs",
			wantKind:   cdb.OptionKindMultiArg,
			wantNumArg: 4,
		},
		{
			name: "single-arg option (Args(1))",
			record: optRecord{
				name:  "onearg",
				attrs: "Args(1) Separate",
			},
			wantSpell:  "-onearg",
			wantKind:   cdb.OptionKindSeparate,
			wantNumArg: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPattern, _ := parseOptRecord(tt.record)
			assert.Equal(t, tt.wantSpell, gotPattern.Spelling)
			assert.Equal(t, tt.wantKind, gotPattern.Kind)
			assert.Equal(t, tt.wantNumArg, gotPattern.NumArgs)
		})
	}
}
