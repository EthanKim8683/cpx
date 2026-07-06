package main

import (
	"fmt"
	"slices"

	"github.com/EthanKim8683/cpx/internal/cdb"
	"github.com/go-json-experiment/json"
)

// defRef wraps a JSON reference to a named def.
type defRef struct {
	Def string `json:"def"` // name of the referenced def
}

// def represents a single TableGen def following the Option class defined in
// llvm/include/llvm/Option/OptParser.td.
type def struct {
	// Superclasses is available on all defs; the remaining fields are only
	// meaningful when Option appears in this list.
	Superclasses []string `json:"!superclasses"`
	Prefixes     []string `json:"Prefixes"`
	Name         string   `json:"Name"`
	Kind         defRef   `json:"Kind"`
	NumArgs      int      `json:"NumArgs"`
	Flags        []defRef `json:"Flags"`
}

// dump is the top-level structure of a TableGen JSON dump.
type dump struct {
	TablegenJSONVersion int                 `json:"!tablegen_json_version"`
	Instanceof          map[string][]string `json:"!instanceof"`
	// Defs captures all non-reserved keys as defs. Requires json/v2 for embed support.
	Defs map[string]def `json:",embed"`
}

// translateDef decomposes a single def into CDB option patterns.
// Only defs inheriting from "Option" are considered.
func translateDef(def def) []cdb.OptionPattern {
	// Only Option defs are relevant to the driver.
	if !slices.Contains(def.Superclasses, "Option") {
		return nil
	}
	// NoDriverOption flags are internal and not exposed to the driver command line.
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
		// KIND_JOINED accepts empty suffix → KIND_JOINED decomposes into Joined and Flag.
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindJoined,
		})
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindFlag,
		})
	case "KIND_SEPARATE":
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindSeparate,
		})
	case "KIND_COMMAJOINED":
		// KIND_COMMAJOINED behaves like KIND_JOINED and consequently decomposes like KIND_JOINED.
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindJoined,
		})
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindFlag,
		})
	case "KIND_MULTIARG":
		partials = append(partials, cdb.OptionPattern{
			Kind:    cdb.OptionKindMultiArg,
			NumArgs: def.NumArgs,
		})
	case "KIND_JOINED_OR_SEPARATE":
		// KIND_JOINED_OR_SEPARATE decomposes into Joined and Separate.
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindJoined,
		})
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindSeparate,
		})
	case "KIND_JOINED_AND_SEPARATE":
		// KIND_JOINED part decomposes like KIND_JOINED.
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindJoinedAndSeparate,
		})
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindSeparate,
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

	// Expand each prefix × kind into a separate pattern.
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
