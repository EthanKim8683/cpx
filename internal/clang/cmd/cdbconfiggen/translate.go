// translate.go translates parsed TableGen JSON dumps into CDB parser configs.

package main

import (
	"fmt"
	"slices"

	"github.com/EthanKim8683/cpx/internal/cdb"
	"github.com/go-json-experiment/json"
)

// defRef is a reference to a def object in the JSON dump.
// See https://llvm.org/docs/TableGen/BackEnds.html#json-reference.
type defRef struct {
	Def string `json:"def"`
}

// def represents a single TableGen def following the Option class defined in
// llvm/include/llvm/Option/OptParser.td. It captures the spelling, kind, flags,
// and prefix information needed to generate CDB option patterns.
type def struct {
	Superclasses []string `json:"!superclasses"`
	Prefixes     []string `json:"Prefixes"`
	Name         string   `json:"Name"`
	Kind         defRef   `json:"Kind"`
	NumArgs      int      `json:"NumArgs"`
	Flags        []defRef `json:"Flags"`
}

// dump represents the full TableGen JSON dump.
// See https://llvm.org/docs/TableGen/BackEnds.html#json-reference.
type dump struct {
	TablegenJSONVersion int                 `json:"!tablegen_json_version"`
	Instanceof          map[string][]string `json:"!instanceof"`
	Defs                map[string]def      `json:",embed"` // all other fields are defs
}

// translateDef decomposes a single def into CDB option patterns.
// Only defs inheriting from "Option" are considered.
func translateDef(def def) []cdb.OptionPattern {
	if !slices.Contains(def.Superclasses, "Option") {
		return nil
	}
	for _, flag := range def.Flags {
		if flag.Def == "NoDriverOption" {
			return nil
		}
	}

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

// unmarshalDump parses a TableGen JSON dump into a dump struct.
func unmarshalDump(data []byte) (*dump, error) {
	var dump dump
	if err := json.Unmarshal(data, &dump); err != nil {
		return nil, fmt.Errorf("unmarshaling dump: %w", err)
	}
	return &dump, nil
}

// translateDump translates an entire TableGen JSON dump into a CDB config.
func translateDump(dump *dump) (*cdb.Config, error) {
	if version := dump.TablegenJSONVersion; version != 1 {
		return nil, fmt.Errorf("unexpected TableGen JSON version: %d", version)
	}

	var patterns []cdb.OptionPattern
	for _, def := range dump.Defs {
		patterns = append(patterns, translateDef(def)...)
	}
	return cdb.NewConfig(patterns), nil
}
