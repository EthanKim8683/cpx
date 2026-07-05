// build.go builds CDB parser configs from parsed TableGen JSON dumps.

package main

import (
	"fmt"
	"slices"

	"github.com/EthanKim8683/cpx/internal/cdb"
	"github.com/go-json-experiment/json"
)

// reference to a def object
// see: https://llvm.org/docs/TableGen/BackEnds.html#json-reference
type defRef struct {
	Def string `json:"def"`
}

// based on the Option class in OptParser.td: https://github.com/llvm/llvm-project/blob/release/22.x/llvm/include/llvm/Option/OptParser.td
// follows translated JSON type from docs: https://llvm.org/docs/TableGen/BackEnds.html#json-reference
type def struct {
	// general to all defs
	Superclasses []string `json:"!superclasses"`

	// Option class-specific fields
	Prefixes []string `json:"Prefixes"`
	Name     string   `json:"Name"`
	Kind     defRef   `json:"Kind"` // Kind -> ref to Kind def
	NumArgs  int      `json:"NumArgs"`
	Flags    []defRef `json:"Flags"` // list<Flag> -> list of refs to Flag def
}

// full json dump type: https://llvm.org/docs/TableGen/BackEnds.html#json-reference
type dump struct {
	TablegenJSONVersion int                 `json:"!tablegen_json_version"`
	Instanceof          map[string][]string `json:"!instanceof"` // to exclude from Defs

	// all other fields are defs
	// JSON embedding is supported by encoding/json/v2: github.com/go-json-experiment/json
	Defs map[string]def `json:",embed"`
}

func buildPatterns(def def) []cdb.OptionPattern {
	// def must inherit from Option to be considered an option
	if !slices.Contains(def.Superclasses, "Option") {
		return nil
	}
	for _, flag := range def.Flags {
		if flag.Def == "NoDriverOption" {
			return nil
		}
	}

	// https://github.com/llvm/llvm-project/blob/release/22.x/llvm/include/llvm/Option/OptParser.td
	var partials []cdb.OptionPattern
	switch def.Kind.Def {
	case "KIND_FLAG":
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindFlag,
		})
	case "KIND_JOINED":
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindJoined,
		})
	case "KIND_SEPARATE":
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindSeparate,
		})
	case "KIND_COMMAJOINED":
		// CommaJoined behaves like Joined with a comma-separated list argument.
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindJoined,
		})
	case "KIND_MULTIARG":
		partials = append(partials, cdb.OptionPattern{
			Kind:    cdb.OptionKindMultiArg,
			NumArgs: def.NumArgs,
		})
	case "KIND_JOINED_OR_SEPARATE":
		// JoinedOrSeparate can be decomposed into Joined and Separate.
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindJoined,
		})
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindSeparate,
		})
	case "KIND_JOINED_AND_SEPARATE":
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindJoinedAndSeparate,
		})
	case "KIND_REMAINING_ARGS":
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindRemainingArgs,
		})
	case "KIND_REMAINING_ARGS_JOINED":
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindRemainingArgsJoined,
		})
	}

	var patterns []cdb.OptionPattern
	for _, prefix := range def.Prefixes {
		for _, partial := range partials {
			patterns = append(patterns, cdb.OptionPattern{
				Spelling: prefix + def.Name,
				Kind:     partial.Kind,
				NumArgs:  partial.NumArgs,
			})
		}
	}
	return patterns
}

func unmarshalDump(data []byte) (*dump, error) {
	var dump dump
	if err := json.Unmarshal(data, &dump); err != nil {
		return nil, fmt.Errorf("unmarshaling dump: %w", err)
	}
	return &dump, nil
}

func buildConfig(dump *dump) (*cdb.Config, error) {
	if version := dump.TablegenJSONVersion; version != 1 {
		return nil, fmt.Errorf("unexpected TableGen JSON version: %d", version)
	}

	var patterns []cdb.OptionPattern
	for _, def := range dump.Defs {
		patterns = append(patterns, buildPatterns(def)...)
	}
	return cdb.NewConfig(patterns), nil
}
