package main

import (
	"testing"

	"github.com/EthanKim8683/cpx/internal/cdb"
	"github.com/stretchr/testify/assert"
)

func TestBuildOptionPatterns(t *testing.T) {
	tests := []struct {
		name   string
		record parsedOptRecord
		want   []cdb.OptionPattern
	}{
		{
			name: "flag option with implicit negation",
			record: parsedOptRecord{
				name:   "fcommon",
			},
			want: []cdb.OptionPattern{
				{
					Spelling: "-fcommon",
					Kind:     cdb.OptionKindFlag,
				},
				{
					Spelling: "-fno-common",
					Kind:     cdb.OptionKindFlag,
				},
			},
		},
		{
			name: "flag option rejecting negative",
			record: parsedOptRecord{
				name:           "Wextra",
				rejectNegative: true,
			},
			want: []cdb.OptionPattern{
				{
					Spelling: "-Wextra",
					Kind:     cdb.OptionKindFlag,
				},
			},
		},
		{
			name: "joined option",
			record: parsedOptRecord{
				name:   "std=",
				joined: true,
			},
			want: []cdb.OptionPattern{
				{
					Spelling: "-std=",
					Kind:     cdb.OptionKindJoined,
				},
			},
		},
		{
			name: "separate option",
			record: parsedOptRecord{
				name:     "o",
				separate: true,
			},
			want: []cdb.OptionPattern{
				{
					Spelling: "-o",
					Kind:     cdb.OptionKindSeparate,
				},
			},
		},
		{
			name: "joined or separate option",
			record: parsedOptRecord{
				name:     "I",
				joined:   true,
				separate: true,
			},
			want: []cdb.OptionPattern{
				{
					Spelling: "-I",
					Kind:     cdb.OptionKindJoinedOrSeparate,
				},
			},
		},
		{
			name: "joined or missing option",
			record: parsedOptRecord{
				name:            "D",
				joinedOrMissing: true,
			},
			want: []cdb.OptionPattern{
				{
					Spelling: "-D",
					Kind:     cdb.OptionKindJoinedOrMissing,
				},
			},
		},
		{
			name: "no driver arg separate maps to flag",
			record: parsedOptRecord{
				name:        "Q",
				separate:    true,
				noDriverArg: true,
			},
			want: []cdb.OptionPattern{
				{
					Spelling: "-Q",
					Kind:     cdb.OptionKindFlag,
				},
			},
		},
		{
			name: "multi-arg option (Args(2))",
			record: parsedOptRecord{
				name:     "sectcreate",
				separate: true,
				args:     2,
			},
			want: []cdb.OptionPattern{
				{
					Spelling: "-sectcreate",
					Kind:     cdb.OptionKindMultiArg,
					NumArgs:  2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildOptionPatterns(tt.record)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestBuildConfig(t *testing.T) {
	records := []parsedOptRecord{
		{
			name:   "fcommon",
		},
		{
			name:           "Wextra",
			rejectNegative: true,
		},
	}

	config := buildConfig(records)

	expectedPrefixes := map[string][]cdb.OptionPattern{
		"-fcommon": {
			{Spelling: "-fcommon", Kind: cdb.OptionKindFlag},
		},
		"-fno-common": {
			{Spelling: "-fno-common", Kind: cdb.OptionKindFlag},
		},
		"-Wextra": {
			{Spelling: "-Wextra", Kind: cdb.OptionKindFlag},
		},
	}

	assert.Equal(t, expectedPrefixes, config.ByPrefix)
}
