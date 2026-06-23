package metadata_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/EthanKim8683/cpx/internal/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRelPath(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		path string
		err  error
	}{
		"ok": {
			path: "main.cpp",
			err:  nil,
		},
		"absolute path": {
			path: "/main.cpp",
			err:  errors.New("path is absolute: /main.cpp"),
		},
		"relative path escapes root directory": {
			path: "../main.cpp",
			err:  errors.New("path escapes root directory: ../main.cpp"),
		},
		"root directory": {
			path: ".",
			// revive:disable-next-line:error-strings
			err: errors.New("path is root directory: ."),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			RelPath, err := metadata.NewRelPath(test.path)
			if test.err != nil {
				require.EqualError(t, err, test.err.Error())
				assert.Empty(t, RelPath)
			} else {
				require.NoError(t, err)
				assert.Equal(t, filepath.Clean(test.path), string(RelPath))
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		args     []string
		metadata *metadata.Metadata
		err      error
	}{
		"g++": {
			args: []string{"g++", "main.cpp"},
			metadata: &metadata.Metadata{
				RelPath: "main.cpp",
				Type:    metadata.MetadataTypeGXX,
				GPP:     &metadata.GPPMetadata{},
			},
		},
		"no arguments": {
			args: []string{},
			err:  errors.New("no arguments provided"),
		},
		"unexpected command": {
			args: []string{"bogus"},
			err:  errors.New("unexpected command: bogus"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			metadata, err := metadata.New(test.args)
			if test.err != nil {
				require.EqualError(t, err, test.err.Error())
				assert.Empty(t, metadata)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.metadata, metadata)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		lhs, rhs *metadata.Metadata
		metadata *metadata.Metadata
		err      error
	}{
		"nil lhs": {
			rhs:      &metadata.Metadata{},
			metadata: &metadata.Metadata{},
		},
		"nil rhs": {
			lhs:      &metadata.Metadata{},
			metadata: &metadata.Metadata{},
		},
		"g++": {
			lhs: &metadata.Metadata{
				RelPath: "main.cpp",
				Type:    metadata.MetadataTypeGXX,
				GPP:     &metadata.GPPMetadata{},
			},
			rhs: &metadata.Metadata{
				RelPath: "main.cpp",
				Type:    metadata.MetadataTypeGXX,
				GPP:     &metadata.GPPMetadata{},
			},
			metadata: &metadata.Metadata{
				RelPath: "main.cpp",
				Type:    metadata.MetadataTypeGXX,
				GPP:     &metadata.GPPMetadata{},
			},
		},
		"different paths": {
			lhs: &metadata.Metadata{
				RelPath: "main1.cpp",
			},
			rhs: &metadata.Metadata{
				RelPath: "main2.cpp",
			},
			err: errors.New("paths do not match: main1.cpp != main2.cpp"),
		},
		"different types": {
			lhs: &metadata.Metadata{
				Type: metadata.MetadataTypeGXX,
				GPP:  &metadata.GPPMetadata{},
			},
			rhs: &metadata.Metadata{
				Type: metadata.MetadataTypeUnspecified,
			},
			err: errors.New("metadata types do not match: g++ != (unspecified)"),
		},
		"unexpected metadata type": {
			lhs: &metadata.Metadata{},
			rhs: &metadata.Metadata{},
			err: errors.New("unexpected metadata type: (unspecified)"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			metadata, err := metadata.Join(test.lhs, test.rhs)
			if test.err != nil {
				require.EqualError(t, err, test.err.Error())
				assert.Empty(t, metadata)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.metadata, metadata)
			}
		})
	}
}
