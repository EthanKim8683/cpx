package cdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    []OptionPattern
		wantSort []string
		wantBC   map[int]string // index → back-chain spelling (omitted = nil)
	}{
		{
			name: "sorts patterns by spelling",
			input: []OptionPattern{
				{Spelling: "-std=", Kind: OptionKindJoined},
				{Spelling: "-o", Kind: OptionKindSeparate},
				{Spelling: "-c", Kind: OptionKindFlag},
			},
			wantSort: []string{"-c", "-o", "-std="},
			wantBC:   map[int]string{},
		},
		{
			name: "back-chain points to longest joined prefix",
			input: []OptionPattern{
				{Spelling: "-std=", Kind: OptionKindJoined},
				{Spelling: "-std=c++", Kind: OptionKindJoined},
				{Spelling: "-std=c++17", Kind: OptionKindJoined},
			},
			wantSort: []string{"-std=", "-std=c++", "-std=c++17"},
			wantBC: map[int]string{
				1: "-std=",
				2: "-std=c++",
			},
		},
		{
			name: "no joined prefix means nil back-chain",
			input: []OptionPattern{
				{Spelling: "-std=", Kind: OptionKindJoined},
				{Spelling: "-x", Kind: OptionKindJoinedAndSeparate},
			},
			wantSort: []string{"-std=", "-x"},
			wantBC:   map[int]string{},
		},
		{
			name: "non-joined kinds have no back-chain",
			input: []OptionPattern{
				{Spelling: "-c", Kind: OptionKindFlag},
				{Spelling: "-o", Kind: OptionKindSeparate},
			},
			wantSort: []string{"-c", "-o"},
			wantBC:   map[int]string{},
		},
		{
			name: "skips non-joined patterns when scanning for joined prefix",
			input: []OptionPattern{
				{Spelling: "-std", Kind: OptionKindFlag},
				{Spelling: "-std=", Kind: OptionKindJoined},
				{Spelling: "-std=c++17", Kind: OptionKindJoined},
			},
			wantSort: []string{"-std", "-std=", "-std=c++17"},
			wantBC: map[int]string{
				2: "-std=",
			},
		},
		{
			name:     "empty input",
			input:    []OptionPattern{},
			wantSort: []string{},
			wantBC:   map[int]string{},
		},
		{
			name: "single non-joined pattern",
			input: []OptionPattern{
				{Spelling: "-c", Kind: OptionKindFlag},
			},
			wantSort: []string{"-c"},
			wantBC:   map[int]string{},
		},
		{
			name: "single joined pattern",
			input: []OptionPattern{
				{Spelling: "-std=", Kind: OptionKindJoined},
			},
			wantSort: []string{"-std="},
			wantBC:   map[int]string{},
		},
		{
			name: "already sorted input",
			input: []OptionPattern{
				{Spelling: "-c", Kind: OptionKindFlag},
				{Spelling: "-o", Kind: OptionKindSeparate},
				{Spelling: "-std=", Kind: OptionKindJoined},
			},
			wantSort: []string{"-c", "-o", "-std="},
			wantBC:   map[int]string{},
		},
		{
			name: "back-chain chains through intermediate joined prefixes",
			input: []OptionPattern{
				{Spelling: "-std=", Kind: OptionKindJoined},
				{Spelling: "-std=c++", Kind: OptionKindJoined},
				{Spelling: "-std=c++17", Kind: OptionKindJoined},
				{Spelling: "-std=c++20", Kind: OptionKindJoined},
			},
			wantSort: []string{"-std=", "-std=c++", "-std=c++17", "-std=c++20"},
			wantBC: map[int]string{
				1: "-std=",
				2: "-std=c++",
				3: "-std=c++",
			},
		},
		{
			name: "non-joined prefix is skipped when scanning for joined prefix",
			input: []OptionPattern{
				{Spelling: "-W", Kind: OptionKindFlag},
				{Spelling: "-Werror", Kind: OptionKindFlag},
				{Spelling: "-Werror=", Kind: OptionKindJoined},
				{Spelling: "-Werror=foo", Kind: OptionKindJoined},
			},
			wantSort: []string{"-W", "-Werror", "-Werror=", "-Werror=foo"},
			wantBC: map[int]string{
				3: "-Werror=",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputCopy := make([]OptionPattern, len(tt.input))
			copy(inputCopy, tt.input)

			cfg := NewConfig(tt.input)

			assert.Equal(t, inputCopy, tt.input, "input must not be mutated")

			require.Len(t, cfg.Patterns, len(tt.wantSort))
			for i, want := range tt.wantSort {
				assert.Equal(t, want, cfg.Patterns[i].Spelling)
			}

			require.Len(t, cfg.BackChains, len(tt.wantSort))
			for i := range cfg.Patterns {
				want, ok := tt.wantBC[i]
				if !ok {
					assert.Nil(t, cfg.BackChains[i])
				} else {
					require.NotNil(t, cfg.BackChains[i])
					assert.Equal(t, want, cfg.BackChains[i].Spelling)
				}
			}
		})
	}
}

func TestIsJoined(t *testing.T) {
	tests := []struct {
		name string
		kind OptionKind
		want bool
	}{
		{name: "Joined", kind: OptionKindJoined, want: true},
		{name: "JoinedAndSeparate", kind: OptionKindJoinedAndSeparate, want: true},
		{name: "RemainingArgsJoined", kind: OptionKindRemainingArgsJoined, want: true},
		{name: "Flag", kind: OptionKindFlag, want: false},
		{name: "Separate", kind: OptionKindSeparate, want: false},
		{name: "MultiArg", kind: OptionKindMultiArg, want: false},
		{name: "RemainingArgs", kind: OptionKindRemainingArgs, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.kind.IsJoined())
		})
	}
}
