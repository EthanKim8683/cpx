package metadata

import (
	"errors"
	"strings"
)

type GXXMetadata struct {
	Command      string   `json:"command"`
	CompileFlags []string `json:"compile_flags"`
	CompileArgs  []string `json:"compile_args"`
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
