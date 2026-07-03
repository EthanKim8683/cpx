// Package cdb provides compiler-agnostic configuration structures, types, and mappings
// for compiling databases.
//
// This package adopts a stateless option parsing strategy inspired by LLVM/Clang's driver architecture
// (https://clang.llvm.org/docs/DriverInternals.html). Rather than resolving complex circular alias
// chains, option negations, and mutual exclusions at generation or parse time (such as GCC's Negative(other)
// attribute system), cdb performs a flat, prefix-based pattern matching pass to slice command-line arguments
// into discrete options and consume the associated number of arguments.
//
// Dynamic option behaviors—including resolving overridden arguments, negation flags, and mutual exclusions—are
// deferred entirely to access or query time. Consumers query the resulting parsed argument list (analogous to
// Clang's InputArgList) using accessors like getLastArg or hasFlag to dynamically resolve the final compiler state.
package cdb

//go:generate go run ./cmd/gccconfiggen -path config/gcc.go

// OptionKind defines the parsing behavior for a compiler option pattern,
// determining how subsequent command-line arguments are consumed.
type OptionKind string

const (
	// OptionKindFlag represents a boolean flag option with no arguments.
	OptionKindFlag OptionKind = "Flag"
	// OptionKindJoined represents an option whose argument is joined directly to the option prefix (e.g. -Ipath).
	OptionKindJoined OptionKind = "Joined"
	// OptionKindSeparate represents an option whose argument is a separate subsequent command-line token (e.g. -o file).
	OptionKindSeparate OptionKind = "Separate"
	// OptionKindJoinedOrSeparate represents an option whose argument can be joined or separate.
	OptionKindJoinedOrSeparate OptionKind = "JoinedOrSeparate"
	// OptionKindCommaJoined represents an option whose arguments are comma-separated within the same token.
	OptionKindCommaJoined OptionKind = "CommaJoined"
	// OptionKindMultiArg represents an option that accepts multiple separate command-line arguments.
	OptionKindMultiArg OptionKind = "MultiArg"
	// OptionKindJoinedAndSeparate represents an option with both a joined and a separate argument.
	OptionKindJoinedAndSeparate OptionKind = "JoinedAndSeparate"
	// OptionKindJoinedOrMissing represents an option with an optional joined argument.
	OptionKindJoinedOrMissing OptionKind = "JoinedOrMissing"
)

// OptionPattern represents a single command-line spelling variant of an option.
type OptionPattern struct {
	// Spelling is the option prefix or full flag string (e.g. "-std=" or "-o").
	Spelling string
	// Kind specifies how the option and its trailing arguments are parsed.
	Kind OptionKind
	// NumArgs is the number of subsequent trailing arguments to consume (only applicable for OptionKindMultiArg).
	NumArgs int
}

// Config represents a complete compiler configuration containing a registry of option patterns.
type Config struct {
	// ByPrefix maps spelling prefixes directly to option pattern specifications.
	ByPrefix map[string][]OptionPattern
}

// NewConfig constructs a Config by indexing the provided option patterns by their spelling prefix.
func NewConfig(patterns []OptionPattern) Config {
	byPrefix := make(map[string][]OptionPattern, len(patterns))
	for _, pattern := range patterns {
		prefix := pattern.Spelling
		byPrefix[prefix] = append(byPrefix[prefix], pattern)
	}
	return Config{ByPrefix: byPrefix}
}
