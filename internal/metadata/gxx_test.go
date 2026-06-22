package metadata_test

import (
	"errors"
	"testing"

	"github.com/EthanKim8683/cpx/internal/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGXXNew(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		args     []string
		metadata *metadata.Metadata
		err      error
	}{
		"attached flags": {
			args: []string{"g++", "main.cpp", "-o", "main", "-std=c++17", "-I./include1", "-I./include2"},
			metadata: &metadata.Metadata{
				RelPath: "main.cpp",
				Type:    metadata.MetadataTypeGXX,
				GXX: &metadata.GXXMetadata{
					Standard:     metadata.GXXStandardCXX17,
					IncludePaths: []string{"./include1", "./include2"},
				},
			},
		},
		"split flags": {
			args: []string{"g++", "main.cpp", "-o", "main", "-std", "c++17", "-I", "./include1", "-I", "./include2"},
			metadata: &metadata.Metadata{
				RelPath: "main.cpp",
				Type:    metadata.MetadataTypeGXX,
				GXX: &metadata.GXXMetadata{
					Standard:     metadata.GXXStandardCXX17,
					IncludePaths: []string{"./include1", "./include2"},
				},
			},
		},
		"no source files": {
			args: []string{"g++"},
			err:  errors.New("no source files specified"),
		},
		"multiple source files": {
			args: []string{"g++", "main1.cpp", "main2.cpp"},
			err:  errors.New("multiple source files specified: [main1.cpp main2.cpp]"),
		},
		"no standard": {
			args: []string{"g++", "main.cpp"},
			metadata: &metadata.Metadata{
				RelPath: "main.cpp",
				Type:    metadata.MetadataTypeGXX,
				GXX: &metadata.GXXMetadata{
					Standard: metadata.GXXStandardDefault,
				},
			},
		},
		"multiple standards": {
			args: []string{"g++", "main.cpp", "-std=c++17", "-std=c++20"},
			metadata: &metadata.Metadata{
				RelPath: "main.cpp",
				Type:    metadata.MetadataTypeGXX,
				GXX: &metadata.GXXMetadata{
					Standard: metadata.GXXStandardCXX20,
				},
			},
		},
		"unexpected standard": {
			args: []string{"g++", "main.cpp", "-std=bogus"},
			err:  errors.New("unexpected standard: bogus"),
		},
		"no standard specified": {
			args: []string{"g++", "main.cpp", "-std"},
			err:  errors.New("no standard specified"),
		},
		"no include path specified": {
			args: []string{"g++", "main.cpp", "-I"},
			err:  errors.New("no include path specified"),
		},
		"no output file specified": {
			args: []string{"g++", "main.cpp", "-o"},
			err:  errors.New("no output file specified"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m, err := metadata.New(test.args)
			if test.err != nil {
				require.EqualError(t, err, test.err.Error())
				assert.Empty(t, m)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.metadata, m)
			}
		})
	}
}

func TestGXXJoin(t *testing.T) {
	t.Parallel()

	lhs := &metadata.Metadata{
		RelPath: "main.cpp",
		Type:    metadata.MetadataTypeGXX,
		GXX: &metadata.GXXMetadata{
			Standard:     metadata.GXXStandardCXX17,
			IncludePaths: []string{"./include1", "./include2"},
		},
	}
	rhs := &metadata.Metadata{
		RelPath: "main.cpp",
		Type:    metadata.MetadataTypeGXX,
		GXX: &metadata.GXXMetadata{
			Standard:     metadata.GXXStandardCXX20,
			IncludePaths: []string{"./include3", "./include4"},
		},
	}
	m, err := metadata.Join(lhs, rhs)
	require.NoError(t, err)
	assert.Equal(t, rhs, m)
}
