package metadata

import (
	"errors"
	"fmt"
	"strings"
)

type GXXStandard string

const (
	GXXStandardDefault GXXStandard = ""
	GXXStandardCXX17   GXXStandard = "c++17"
	GXXStandardCXX20   GXXStandard = "c++20"
	GXXStandardCXX23   GXXStandard = "c++23"
)

func (s GXXStandard) String() string {
	if s == GXXStandardDefault {
		return "(default)"
	}

	return string(s)
}

func NewGXXStandard(standard string) (GXXStandard, error) {
	switch standard {
	case "", "c++17", "c++20", "c++23":
		return GXXStandard(standard), nil
	default:
		return "", fmt.Errorf("unexpected standard: %s", standard)
	}
}

type GXXMetadata struct {
	Standard     GXXStandard `json:"standard"`
	IncludePaths []string    `json:"include_paths"`
}

type gxxMetadataBuilder struct {
	errs         error
	sourcePaths  []RelPath
	standards    []GXXStandard
	includePaths []string
}

func (b *gxxMetadataBuilder) addError(err error) {
	b.errs = errors.Join(b.errs, err)
}

func (b *gxxMetadataBuilder) addSourcePath(arg string) {
	sourcePath, err := NewRelPath(arg)
	if err != nil {
		b.addError(err)
		return
	}
	b.sourcePaths = append(b.sourcePaths, sourcePath)
}

func (b *gxxMetadataBuilder) addStandard(arg string) {
	standard, err := NewGXXStandard(arg)
	if err != nil {
		b.addError(err)
		return
	}
	b.standards = append(b.standards, standard)
}

func (b *gxxMetadataBuilder) addIncludePath(arg string) {
	b.includePaths = append(b.includePaths, arg)
}

func (b *gxxMetadataBuilder) build() (*Metadata, error) {
	var sourcePath RelPath
	switch {
	case len(b.sourcePaths) == 0:
		b.addError(errors.New("no source files specified"))
	case len(b.sourcePaths) > 1:
		b.addError(fmt.Errorf("multiple source files specified: %v", b.sourcePaths))
	default:
		sourcePath = b.sourcePaths[0]
	}

	var standard GXXStandard
	switch {
	case len(b.standards) == 0:
		standard = GXXStandardDefault
	default:
		standard = b.standards[len(b.standards)-1]
	}

	if b.errs != nil {
		return nil, b.errs
	}

	return &Metadata{
		RelPath: sourcePath,
		Type:    MetadataTypeGXX,
		GXX: &GXXMetadata{
			Standard:     standard,
			IncludePaths: b.includePaths,
		},
	}, nil
}

func gxxNew(args []string) (*Metadata, error) {
	var errs error
	var b gxxMetadataBuilder

	i := 0
	eatArg := func() (string, bool) {
		if i >= len(args) {
			return "", false
		}

		arg := args[i]
		i++

		return arg, true
	}

	for {
		arg, ok := eatArg()
		if !ok {
			break
		}

		if strings.HasPrefix(arg, "-") {
			switch {
			case arg == "-o":
				_, ok := eatArg()
				if !ok {
					errs = errors.Join(errs, errors.New("no output file specified"))
					break
				}
			case arg == "-std":
				standardArg, ok := eatArg()
				if !ok {
					errs = errors.Join(errs, errors.New("no standard specified"))
					break
				}
				b.addStandard(standardArg)
			case arg == "-I":
				includePathArg, ok := eatArg()
				if !ok {
					errs = errors.Join(errs, errors.New("no include path specified"))
					break
				}
				b.addIncludePath(includePathArg)
			case strings.HasPrefix(arg, "-std="):
				b.addStandard(strings.TrimPrefix(arg, "-std="))
			case strings.HasPrefix(arg, "-I"):
				b.addIncludePath(strings.TrimPrefix(arg, "-I"))
			}
		} else {
			b.addSourcePath(arg)
		}
	}

	m, err := b.build()
	errs = errors.Join(errs, err)

	if errs != nil {
		return nil, errs
	}
	return m, nil
}

func gxxJoin(_, rhs *Metadata) (*Metadata, error) {
	return rhs, nil
}
