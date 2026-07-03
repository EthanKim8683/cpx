// parse.go parses GCC option record properties to extract compilation database option traits.
// For option property specifications, see the GCC Option File format specs
// (https://gcc.gnu.org/onlinedocs/gccint/Option-file-format.html).

package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// parenRE matches parenthesis blocks (e.g. `(flag)`) to strip them from property strings.
var parenRE = regexp.MustCompile(`\([^)]+\)`)

// parsedOptRecord represents the parsed and decoded properties of a GCC option record.
type parsedOptRecord struct {
	name            string
	rejectDriver    bool
	rejectNegative  bool
	joined          bool
	separate        bool
	joinedOrMissing bool
	args            int
	noDriverArg     bool
}

// hasProp checks if a property word is present in the properties string,
// stripping any parenthesized arguments first to avoid collisions with nested parameters.
func hasProp(prop, props string) bool {
	// Stripping parenthesized sub-blocks prevents matching nested option names or arguments
	// (e.g. matching a nested parameter name instead of the parent attribute word).
	props = parenRE.ReplaceAllString(props, "")

	// Enforces exact word-boundary matching to mimic flag_set_p in opt-functions.awk
	// (https://github.com/gcc-mirror/gcc/blob/releases/gcc-16/gcc/opt-functions.awk).
	return strings.Contains(" "+props+" ", " "+prop+" ")
}

// propArgs extracts the arguments of a parameterized property (e.g. Args(N) or Var(name)).
func propArgs(name, props string) string {
	_, s, found := strings.Cut(" "+props, " "+name+"(")
	if !found {
		return ""
	}
	s, found = strings.CutPrefix(s, "{")
	if found {
		// Braced lists are unwrapped to handle comma-containing parameters grouped by GCC.
		// See opt_args in opt-functions.awk (https://github.com/gcc-mirror/gcc/blob/releases/gcc-16/gcc/opt-functions.awk).
		s, _, _ = strings.Cut(s, "})")
	} else {
		s, _, _ = strings.Cut(s, ")")
	}
	return s
}

// parseOptRecord decodes the properties of a raw option record into a parsedOptRecord.
func parseOptRecord(record optRecord) (parsedOptRecord, error) {
	var parsed parsedOptRecord
	parsed.name = record.name
	props := record.props
	parsed.rejectDriver = hasProp("RejectDriver", props)
	parsed.rejectNegative = hasProp("RejectNegative", props)
	parsed.joined = hasProp("Joined", props)
	parsed.separate = hasProp("Separate", props)
	parsed.joinedOrMissing = hasProp("JoinedOrMissing", props)
	if s := propArgs("Args", props); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			return parsed, fmt.Errorf("parsing Args(n): atoi(n): %w", err)
		}
		parsed.args = n
	}
	parsed.noDriverArg = hasProp("NoDriverArg", props)
	return parsed, nil
}
