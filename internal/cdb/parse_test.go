package cdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindPattern(t *testing.T) {
	// Sorted: - (0 RemainingArgs), -MF (1 MultiArg), -c (2 Flag),
	//         -o (3 Separate), -std (4 Flag), -std= (5 Joined),
	//         -std=c++17 (6 Joined)
	//
	// backChains: [5]=nil (no joined prefix for -std=),
	//             [6]=-std= (longest joined prefix for -std=c++17)
	cfg := NewConfig([]OptionPattern{
		{Spelling: "-", Kind: OptionKindRemainingArgs},
		{Spelling: "-MF", Kind: OptionKindMultiArg, NumArgs: 2},
		{Spelling: "-c", Kind: OptionKindFlag},
		{Spelling: "-o", Kind: OptionKindSeparate},
		{Spelling: "-std", Kind: OptionKindFlag},
		{Spelling: "-std=", Kind: OptionKindJoined},
		{Spelling: "-std=c++17", Kind: OptionKindJoined},
	})

	tests := []struct {
		name     string
		arg      string
		wantNil  bool
		wantSp   string
		wantKind OptionKind
	}{
		{
			name:     "exact match on Flag",
			arg:      "-c",
			wantSp:   "-c",
			wantKind: OptionKindFlag,
		},
		{
			name:     "exact match on Separate",
			arg:      "-o",
			wantSp:   "-o",
			wantKind: OptionKindSeparate,
		},
		{
			name:     "exact match on MultiArg",
			arg:      "-MF",
			wantSp:   "-MF",
			wantKind: OptionKindMultiArg,
		},
		{
			name:     "exact match on Flag variant of joined spelling",
			arg:      "-std",
			wantSp:   "-std",
			wantKind: OptionKindFlag,
		},
		{
			name:     "exact match on Joined with no joined prefix",
			arg:      "-std=",
			wantNil: true,
		},
		{
			name:     "exact match on Joined follows back-chain",
			arg:      "-std=c++17",
			wantSp:   "-std=",
			wantKind: OptionKindJoined,
		},
		{
			// -std=c++20 misses on binary search (past -std=c++17).
			// Miss loop: -std=c++17 is not a prefix, follow back-chain
			// to -std= which is. Multi-step back-chain traversal.
			name:     "prefix match for longer joined arg",
			arg:      "-std=c++20",
			wantSp:   "-std=",
			wantKind: OptionKindJoined,
		},
		{
			name:     "prefix match skips non-joined patterns",
			arg:      "-std=gnu++17",
			wantSp:   "-std=",
			wantKind: OptionKindJoined,
		},
		{
			name:    "unknown option",
			arg:     "-unknown",
			wantNil: true,
		},
		{
			name:    "no prefix match",
			arg:     "-z",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findPattern(cfg, tt.arg)
			if tt.wantNil {
				assert.Nil(t, got)
			} else {
				require.NotNil(t, got)
				assert.Equal(t, tt.wantSp, got.Spelling)
				assert.Equal(t, tt.wantKind, got.Kind)
			}
		})
	}
}

func TestParse(t *testing.T) {
	cfg := NewConfig([]OptionPattern{
		{Spelling: "-", Kind: OptionKindRemainingArgs},
		{Spelling: "-MF", Kind: OptionKindMultiArg, NumArgs: 2},
		{Spelling: "-c", Kind: OptionKindFlag},
		{Spelling: "-o", Kind: OptionKindSeparate},
		{Spelling: "-std", Kind: OptionKindFlag},
		{Spelling: "-std=", Kind: OptionKindJoined},
		{Spelling: "-std=c++17", Kind: OptionKindJoined},
	})

	tests := []struct {
		name       string
		argv       []string
		wantErr    bool
		wantName   string
		wantOpts   []Option
		wantArgs   []string
		errContain string
	}{
		{
			name:     "only command name",
			argv:     []string{"cc"},
			wantName: "cc",
		},
		{
			name:     "flag option",
			argv:     []string{"cc", "-c", "main.c"},
			wantName: "cc",
			wantOpts: []Option{
				{Pattern: OptionPattern{Spelling: "-c", Kind: OptionKindFlag}, Args: []string{"-c"}},
			},
			wantArgs: []string{"main.c"},
		},
		{
			name:     "joined option via prefix match",
			argv:     []string{"cc", "-std=c++17", "main.c"},
			wantName: "cc",
			wantOpts: []Option{
				{Pattern: OptionPattern{Spelling: "-std=", Kind: OptionKindJoined}, Args: []string{"-std=c++17"}},
			},
			wantArgs: []string{"main.c"},
		},
		{
			name:     "separate option",
			argv:     []string{"cc", "-o", "out", "main.c"},
			wantName: "cc",
			wantOpts: []Option{
				{Pattern: OptionPattern{Spelling: "-o", Kind: OptionKindSeparate}, Args: []string{"-o", "out"}},
			},
			wantArgs: []string{"main.c"},
		},
		{
			name:     "multi-arg option",
			argv:     []string{"cc", "-MF", "a", "b", "main.c"},
			wantName: "cc",
			wantOpts: []Option{
				{Pattern: OptionPattern{Spelling: "-MF", Kind: OptionKindMultiArg, NumArgs: 2}, Args: []string{"-MF", "a", "b"}},
			},
			wantArgs: []string{"main.c"},
		},
		{
			name:     "remaining args option consumes tail",
			argv:     []string{"cc", "-", "a", "b", "c"},
			wantName: "cc",
			wantOpts: []Option{
				{Pattern: OptionPattern{Spelling: "-", Kind: OptionKindRemainingArgs}, Args: []string{"-", "a", "b", "c"}},
			},
		},
		{
			name:     "unknown options become positional args",
			argv:     []string{"cc", "-foo", "-bar"},
			wantName: "cc",
			wantArgs: []string{"-foo", "-bar"},
		},
		{
			name:     "mixed options and args",
			argv:     []string{"cc", "-c", "-o", "out", "main.c", "-foo"},
			wantName: "cc",
			wantOpts: []Option{
				{Pattern: OptionPattern{Spelling: "-c", Kind: OptionKindFlag}, Args: []string{"-c"}},
				{Pattern: OptionPattern{Spelling: "-o", Kind: OptionKindSeparate}, Args: []string{"-o", "out"}},
			},
			wantArgs: []string{"main.c", "-foo"},
		},
		{
			name:       "separate missing argument",
			argv:       []string{"cc", "-o"},
			wantErr:    true,
			errContain: "option -o takes 1 arguments",
		},
		{
			name:       "multi-arg missing arguments",
			argv:       []string{"cc", "-MF", "a"},
			wantErr:    true,
			errContain: "option -MF takes 2 arguments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(cfg, tt.argv)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContain)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantName, got.Name)
			assert.Equal(t, tt.wantOpts, got.Options)
			assert.Equal(t, tt.wantArgs, got.Args)
		})
	}
}
