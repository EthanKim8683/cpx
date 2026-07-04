// Option patterns building implementation from Clang TableGen JSON.
//
// This file builds option configurations from parsed TableGen JSON records.

package main

import (
	"slices"

	"github.com/EthanKim8683/cpx/internal/cdb"
)

func buildOptionPatterns(record optionDef) []cdb.OptionPattern {
	if !slices.Contains(record.Superclasses, "Option") {
		return nil
	}
	for _, flag := range record.Flags {
		if flag.Def == "NoDriverOption" {
			return nil
		}
	}

	// https://github.com/llvm/llvm-project/blob/release/22.x/llvm/include/llvm/Option/OptParser.td
	var partials []cdb.OptionPattern
	switch record.Kind.Def {
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
		// CommaJoined is Joined with a comma-separated list argument.
		partials = append(partials, cdb.OptionPattern{
			Kind: cdb.OptionKindJoined,
		})
	case "KIND_MULTIARG":
		partials = append(partials, cdb.OptionPattern{
			Kind:    cdb.OptionKindMultiArg,
			NumArgs: record.NumArgs,
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
	for _, prefix := range record.Prefixes {
		for _, partial := range partials {
			patterns = append(patterns, cdb.OptionPattern{
				Spelling: prefix + record.Name,
				Kind:     partial.Kind,
				NumArgs:  partial.NumArgs,
			})
		}
	}
	return patterns
}

func buildConfig(parsedDump *parsedDump) cdb.Config {
	var patterns []cdb.OptionPattern
	for _, record := range parsedDump.Options {
		patterns = append(patterns, buildOptionPatterns(record)...)
	}
	return cdb.NewConfig(patterns)
}
