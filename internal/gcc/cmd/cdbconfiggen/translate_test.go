package main

import (
	"testing"

	"github.com/EthanKim8683/cpx/internal/cdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeOptRecords(t *testing.T) {
	t.Run("merges duplicate names", func(t *testing.T) {
		records := []optRecord{
			{name: "foo", props: "Joined"},
			{name: "bar", props: "Separate"},
			{name: "foo", props: "RejectNegative"},
		}
		got := mergeOptRecords(records)
		require.Len(t, got, 2)
		require.Equal(t, optRecord{name: "bar", props: "Separate"}, got[0])
		require.Equal(t, optRecord{name: "foo", props: "Joined RejectNegative"}, got[1])
	})

	t.Run("no duplicates", func(t *testing.T) {
		records := []optRecord{
			{name: "a", props: "Joined"},
			{name: "b", props: "Separate"},
		}
		got := mergeOptRecords(records)
		require.Len(t, got, 2)
	})
}

func TestHasProp(t *testing.T) {
	tests := []struct {
		name  string
		prop  string
		props string
		want  bool
	}{
		{
			name:  "present",
			prop:  "Joined",
			props: "Common Driver Joined",
			want:  true,
		},
		{
			name:  "only value",
			prop:  "Joined",
			props: "Joined",
			want:  true,
		},
		{
			name:  "absent",
			prop:  "Joined",
			props: "Common Driver Separate",
			want:  false,
		},
		{
			name:  "empty props",
			prop:  "Joined",
			props: "",
			want:  false,
		},
		{
			name:  "inside parens is ignored",
			prop:  "Joined",
			props: "Foo(Joined) Bar",
			want:  false,
		},
		{
			name:  "after paren group",
			prop:  "Joined",
			props: "Foo(Bar) Joined",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasProp(tt.prop, tt.props)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPropArgs(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		props string
		want  string
	}{
		{
			name:  "simple args",
			key:   "Args",
			props: "Separate Args(2)",
			want:  "2",
		},
		{
			name:  "args with other props",
			key:   "Args",
			props: "Joined Separate Args(3) Warning",
			want:  "3",
		},
		{
			name:  "brace-wrapped value",
			key:   "Var",
			props: "Var({foo})",
			want:  "foo",
		},
		{
			name:  "not found",
			key:   "Args",
			props: "Joined Separate",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := propArgs(tt.key, tt.props)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNegative(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "ffoo", in: "ffoo", want: "fno-foo"},
		{name: "Wextra", in: "Wextra", want: "Wno-extra"},
		{name: "msse", in: "msse", want: "mno-sse"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := negative(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTranslateOptRecord(t *testing.T) {
	tests := []struct {
		name   string
		record optRecord
		want   []cdb.OptionPattern
	}{
		{
			name:   "RejectDriver is skipped",
			record: optRecord{name: "foo", props: "Joined RejectDriver"},
		},
		{
			name:   "flag with no properties",
			record: optRecord{name: "static", props: ""},
			want:   []cdb.OptionPattern{{Spelling: "-static", Kind: cdb.OptionKindFlag}},
		},
		{
			name:   "flag with negation",
			record: optRecord{name: "ffoo", props: ""},
			want: []cdb.OptionPattern{
				{Spelling: "-ffoo", Kind: cdb.OptionKindFlag},
				{Spelling: "-fno-foo", Kind: cdb.OptionKindFlag},
			},
		},
		{
			name:   "flag with RejectNegative",
			record: optRecord{name: "ffoo", props: "RejectNegative"},
			want:   []cdb.OptionPattern{{Spelling: "-ffoo", Kind: cdb.OptionKindFlag}},
		},
		{
			name:   "Joined",
			record: optRecord{name: "std=", props: "Joined"},
			want:   []cdb.OptionPattern{{Spelling: "-std=", Kind: cdb.OptionKindJoined}},
		},
		{
			name:   "Separate",
			record: optRecord{name: "o", props: "Separate"},
			want:   []cdb.OptionPattern{{Spelling: "-o", Kind: cdb.OptionKindSeparate}},
		},
		{
			name:   "Separate with NoDriverArg becomes Flag",
			record: optRecord{name: "o", props: "Separate NoDriverArg"},
			want:   []cdb.OptionPattern{{Spelling: "-o", Kind: cdb.OptionKindFlag}},
		},
		{
			name:   "Separate with Args becomes MultiArg",
			record: optRecord{name: "MF", props: "Separate Args(2)"},
			want:   []cdb.OptionPattern{{Spelling: "-MF", Kind: cdb.OptionKindMultiArg, NumArgs: 2}},
		},
		{
			name:   "Joined + Separate produces both",
			record: optRecord{name: "o", props: "Joined Separate"},
			want: []cdb.OptionPattern{
				{Spelling: "-o", Kind: cdb.OptionKindJoined},
				{Spelling: "-o", Kind: cdb.OptionKindSeparate},
			},
		},
		{
			name:   "JoinedOrMissing decomposes into Flag and Joined",
			record: optRecord{name: "x", props: "JoinedOrMissing"},
			want: []cdb.OptionPattern{
				{Spelling: "-x", Kind: cdb.OptionKindFlag},
				{Spelling: "-x", Kind: cdb.OptionKindJoined},
			},
		},
		{
			name:   "non-negatable name has no negation",
			record: optRecord{name: "static", props: ""},
			want:   []cdb.OptionPattern{{Spelling: "-static", Kind: cdb.OptionKindFlag}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := translateOptRecord(tt.record)
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestTranslateOptRecords(t *testing.T) {
	t.Run("normal input", func(t *testing.T) {
		records := []optRecord{
			{name: "foo", props: "Joined"},
			{name: "bar", props: "Separate"},
		}
		got := translateOptRecords(records)
		require.NotNil(t, got)
		require.Len(t, got.ByPrefix["-foo"], 1)
		require.Equal(t, cdb.OptionKindJoined, got.ByPrefix["-foo"][0].Kind)
		require.Len(t, got.ByPrefix["-bar"], 1)
		require.Equal(t, cdb.OptionKindSeparate, got.ByPrefix["-bar"][0].Kind)
	})

	t.Run("nil input", func(t *testing.T) {
		got := translateOptRecords(nil)
		require.NotNil(t, got)
		require.Empty(t, got.ByPrefix)
	})
}
