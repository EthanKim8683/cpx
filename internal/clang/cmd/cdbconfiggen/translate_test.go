package main

import (
	"testing"

	"github.com/EthanKim8683/cpx/internal/cdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTranslateDef(t *testing.T) {
	tests := []struct {
		name string
		def  def
		want []cdb.OptionPattern
	}{
		{
			name: "non-option def is skipped",
			def: def{
				Superclasses: []string{"Base"},
				Prefixes:     []string{"-"},
				Name:         "foo",
				Kind:         defRef{Def: "KIND_FLAG"},
			},
		},
		{
			name: "NoDriverOption flag is skipped",
			def: def{
				Superclasses: []string{"Option"},
				Prefixes:     []string{"-"},
				Name:         "foo",
				Kind:         defRef{Def: "KIND_FLAG"},
				Flags:        []defRef{{Def: "NoDriverOption"}},
			},
		},
		{
			name: "KIND_FLAG",
			def: def{
				Superclasses: []string{"Option"},
				Prefixes:     []string{"-"},
				Name:         "foo",
				Kind:         defRef{Def: "KIND_FLAG"},
			},
			want: []cdb.OptionPattern{
				{Spelling: "-foo", Kind: cdb.OptionKindFlag},
			},
		},
		{
			name: "KIND_JOINED",
			def: def{
				Superclasses: []string{"Option"},
				Prefixes:     []string{"-"},
				Name:         "std=",
				Kind:         defRef{Def: "KIND_JOINED"},
			},
			want: []cdb.OptionPattern{
				{Spelling: "-std=", Kind: cdb.OptionKindJoined},
			},
		},
		{
			name: "KIND_SEPARATE",
			def: def{
				Superclasses: []string{"Option"},
				Prefixes:     []string{"-"},
				Name:         "o",
				Kind:         defRef{Def: "KIND_SEPARATE"},
			},
			want: []cdb.OptionPattern{
				{Spelling: "-o", Kind: cdb.OptionKindSeparate},
			},
		},
		{
			name: "KIND_COMMAJOINED maps to Joined",
			def: def{
				Superclasses: []string{"Option"},
				Prefixes:     []string{"-"},
				Name:         "Wa,",
				Kind:         defRef{Def: "KIND_COMMAJOINED"},
			},
			want: []cdb.OptionPattern{
				{Spelling: "-Wa,", Kind: cdb.OptionKindJoined},
			},
		},
		{
			name: "KIND_MULTIARG with NumArgs",
			def: def{
				Superclasses: []string{"Option"},
				Prefixes:     []string{"-"},
				Name:         "MF",
				Kind:         defRef{Def: "KIND_MULTIARG"},
				NumArgs:      2,
			},
			want: []cdb.OptionPattern{
				{Spelling: "-MF", Kind: cdb.OptionKindMultiArg, NumArgs: 2},
			},
		},
		{
			name: "KIND_JOINED_OR_SEPARATE decomposes into both",
			def: def{
				Superclasses: []string{"Option"},
				Prefixes:     []string{"-"},
				Name:         "o",
				Kind:         defRef{Def: "KIND_JOINED_OR_SEPARATE"},
			},
			want: []cdb.OptionPattern{
				{Spelling: "-o", Kind: cdb.OptionKindJoined},
				{Spelling: "-o", Kind: cdb.OptionKindSeparate},
			},
		},
		{
			name: "KIND_JOINED_AND_SEPARATE",
			def: def{
				Superclasses: []string{"Option"},
				Prefixes:     []string{"-"},
				Name:         "MF",
				Kind:         defRef{Def: "KIND_JOINED_AND_SEPARATE"},
			},
			want: []cdb.OptionPattern{
				{Spelling: "-MF", Kind: cdb.OptionKindJoinedAndSeparate},
			},
		},
		{
			name: "KIND_REMAINING_ARGS",
			def: def{
				Superclasses: []string{"Option"},
				Prefixes:     []string{"-"},
				Name:         "",
				Kind:         defRef{Def: "KIND_REMAINING_ARGS"},
			},
			want: []cdb.OptionPattern{
				{Spelling: "-", Kind: cdb.OptionKindRemainingArgs},
			},
		},
		{
			name: "KIND_REMAINING_ARGS_JOINED",
			def: def{
				Superclasses: []string{"Option"},
				Prefixes:     []string{"--"},
				Name:         "CLASSPATH=",
				Kind:         defRef{Def: "KIND_REMAINING_ARGS_JOINED"},
			},
			want: []cdb.OptionPattern{
				{Spelling: "--CLASSPATH=", Kind: cdb.OptionKindRemainingArgsJoined},
			},
		},
		{
			name: "multiple prefixes expand into separate patterns",
			def: def{
				Superclasses: []string{"Option"},
				Prefixes:     []string{"-", "--"},
				Name:         "version",
				Kind:         defRef{Def: "KIND_FLAG"},
			},
			want: []cdb.OptionPattern{
				{Spelling: "-version", Kind: cdb.OptionKindFlag},
				{Spelling: "--version", Kind: cdb.OptionKindFlag},
			},
		},
		{
			name: "empty prefixes produce no patterns",
			def: def{
				Superclasses: []string{"Option"},
				Prefixes:     []string{},
				Name:         "foo",
				Kind:         defRef{Def: "KIND_FLAG"},
			},
			want: nil,
		},
		{
			name: "unknown kind produces no patterns",
			def: def{
				Superclasses: []string{"Option"},
				Prefixes:     []string{"-"},
				Name:         "foo",
				Kind:         defRef{Def: "KIND_UNKNOWN"},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := translateDef(tt.def)
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUnmarshalDump(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		data := []byte(`{"!tablegen_json_version": 1, "!instanceof": {}, "foo": {"!superclasses": ["Option"], "Prefixes": ["-"], "Name": "foo", "Kind": {"def": "KIND_FLAG"}, "NumArgs": 0, "Flags": []}}`)
		got, err := unmarshalDump(data)
		require.NoError(t, err)
		require.Equal(t, 1, got.TablegenJSONVersion)
		require.Len(t, got.Defs, 1)
		require.Contains(t, got.Defs, "foo")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := unmarshalDump([]byte("not json"))
		require.Error(t, err)
	})
}

func TestTranslateDump(t *testing.T) {
	t.Run("valid dump", func(t *testing.T) {
		d := &dump{
			TablegenJSONVersion: 1,
			Defs: map[string]def{
				"foo": {
					Superclasses: []string{"Option"},
					Prefixes:     []string{"-"},
					Name:         "foo",
					Kind:         defRef{Def: "KIND_FLAG"},
				},
				"bar": {
					Superclasses: []string{"Option"},
					Prefixes:     []string{"-"},
					Name:         "bar",
					Kind:         defRef{Def: "KIND_JOINED"},
				},
				"nonoption": {
					Superclasses: []string{"Base"},
					Prefixes:     []string{"-"},
					Name:         "baz",
					Kind:         defRef{Def: "KIND_FLAG"},
				},
			},
		}
		got, err := translateDump(d)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Len(t, got.ByPrefix["-foo"], 1)
		require.Equal(t, cdb.OptionKindFlag, got.ByPrefix["-foo"][0].Kind)
		require.Len(t, got.ByPrefix["-bar"], 1)
		require.Equal(t, cdb.OptionKindJoined, got.ByPrefix["-bar"][0].Kind)
		require.Empty(t, got.ByPrefix["-baz"])
	})

	t.Run("wrong version", func(t *testing.T) {
		d := &dump{TablegenJSONVersion: 2}
		_, err := translateDump(d)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unexpected TableGen JSON version")
	})
}
