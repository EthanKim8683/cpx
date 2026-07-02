// Package cdb provides compiler-agnostic configuration structures, types, and mappings
// for compiling databases.
package cdb

// Option represents a parsed option instance in a compile command.
type Option struct {
	ID   int
	Args []string
}

// OptionAlias represents a mapping from a spelling pattern to a canonical ID and implicit arguments.
type OptionAlias struct {
	ID        int
	AliasArgs []string
}

// OptionPattern represents a single command-line spelling variant of an option.
type OptionPattern struct {
	Spelling string
	Kind     OptionKind
	NumArgs  int
}

// OptionKind defines the parsing behavior for a compiler option pattern.
type OptionKind string

const (
	// OptionKindFlag represents a boolean flag option with no arguments.
	OptionKindFlag OptionKind = "Flag"
	// OptionKindJoined represents an option whose argument is joined to the flag spelling.
	OptionKindJoined OptionKind = "Joined"
	// OptionKindSeparate represents an option whose argument is a separate subsequent token in argv.
	OptionKindSeparate OptionKind = "Separate"
	// OptionKindJoinedOrSeparate represents an option whose argument can be joined or separate.
	OptionKindJoinedOrSeparate OptionKind = "JoinedOrSeparate"
	// OptionKindCommaJoined represents an option whose arguments are comma-separated.
	OptionKindCommaJoined OptionKind = "CommaJoined"
	// OptionKindMultiArg represents an option that accepts multiple separate arguments.
	OptionKindMultiArg OptionKind = "MultiArg"
	// OptionKindJoinedAndSeparate represents an option with both a joined and separate argument.
	OptionKindJoinedAndSeparate OptionKind = "JoinedAndSeparate"
	// OptionKindJoinedOrMissing represents an option with an optional joined argument.
	OptionKindJoinedOrMissing OptionKind = "JoinedOrMissing"
)

// Config represents a compiled set of option configurations for a specific compiler.
type Config struct {
	Aliases map[OptionPattern]OptionAlias
}
