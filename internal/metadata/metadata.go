package metadata

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type RelPath string

func NewRelPath(path string) (RelPath, error) {
	path = filepath.Clean(path)

	if filepath.IsAbs(path) {
		return "", fmt.Errorf("path is absolute: %s", path)
	}

	if strings.HasPrefix(path, "..") {
		return "", fmt.Errorf("path escapes root directory: %s", path)
	}

	if path == "." {
		return "", fmt.Errorf("path is root directory: %s", path)
	}

	return RelPath(path), nil
}

type MetadataType string

const (
	MetadataTypeUnspecified MetadataType = ""
	MetadataTypeGXX         MetadataType = "g++"
)

func (t MetadataType) String() string {
	if t == MetadataTypeUnspecified {
		return "(unspecified)"
	}

	return string(t)
}

type Metadata struct {
	RelPath RelPath      `json:"path"`
	Type    MetadataType `json:"type"`
	GPP     *GPPMetadata `json:"g++"`
}

func New(args []string) (*Metadata, error) {
	if len(args) == 0 {
		return nil, errors.New("no arguments provided")
	}

	switch args[0] {
	case "g++":
		return gppNew(args[1:])
	default:
		return nil, fmt.Errorf("unexpected command: %s", args[0])
	}
}

func Join(lhs, rhs *Metadata) (*Metadata, error) {
	if lhs == nil {
		return rhs, nil
	}
	if rhs == nil {
		return lhs, nil
	}

	if lhs.RelPath != rhs.RelPath {
		return nil, fmt.Errorf("paths do not match: %s != %s", lhs.RelPath, rhs.RelPath)
	}

	if lhs.Type != rhs.Type {
		return nil, fmt.Errorf("metadata types do not match: %s != %s", lhs.Type, rhs.Type)
	}

	switch lhs.Type {
	case MetadataTypeGXX:
		return gppJoin(lhs, rhs)
	default:
		return nil, fmt.Errorf("unexpected metadata type: %s", lhs.Type)
	}
}
