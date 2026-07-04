// parse.go parses raw TableGen JSON dumps into structured representations.

package main

import (
	"errors"
	"fmt"
	"slices"

	"github.com/go-json-experiment/json"
)

// errUnexpectedTablegenJSONVersion is returned when encountering an unhandled TableGen JSON version.
var errUnexpectedTablegenJSONVersion = errors.New("unexpected TableGen JSON version")

// defRef sparsely mirrors the structure of a reference to a def object:
// https://llvm.org/docs/TableGen/BackEnds.html#json-reference
type defRef struct {
	Def string `json:"def"`
}

// optionDef sparsely mirrors the Option class:
// https://github.com/llvm/llvm-project/blob/release/22.x/llvm/include/llvm/Option/OptParser.td
type optionDef struct {
	// Must contain Option for the below properties to be valid.
	Superclasses []string `json:"!superclasses"`

	Prefixes []string `json:"Prefixes"` // list<string> Prefixes
	Name     string   `json:"Name"`     // string Name
	Kind     defRef   `json:"Kind"`     // OptionKind Kind
	NumArgs  int      `json:"NumArgs"`  // int NumArgs
	Flags    []defRef `json:"Flags"`    // list<OptionFlag> Flags
}

// parsedDump mirrors the structure of a TableGen JSON dump:
// https://llvm.org/docs/TableGen/BackEnds.html#json-reference
type parsedDump struct {
	TablegenJSONVersion int                 `json:"!tablegen_json_version"`
	Instanceof          map[string][]string `json:"!instanceof"`
	// JSON embedding is supported by encoding/json/v2:
	// github.com/go-json-experiment/json
	Options map[string]optionDef `json:",embed"`
}

// parseDump parses a raw TableGen JSON dump into a structured representation.
func parseDump(data []byte) (*parsedDump, error) {
	var d parsedDump
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, fmt.Errorf("unmarshalling dump: %w", err)
	}

	// Currently only version 1 is supported.
	if d.TablegenJSONVersion != 1 {
		return nil, errUnexpectedTablegenJSONVersion
	}

	// All defs were unmarshalled as option defs. Remove non-option defs.
	for k, option := range d.Options {
		if !slices.Contains(option.Superclasses, "Option") {
			delete(d.Options, k)
		}
	}
	return &d, nil
}
