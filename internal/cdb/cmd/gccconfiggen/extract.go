// extract.go extracts individual option record blocks from GCC option file content.
// It handles comment stripping, newline normalization, and excludes metadata definitions.
// For the underlying file layout specifications, see the GCC Option File format specs
// (https://gcc.gnu.org/onlinedocs/gccint/Option-file-format.html).

package main

import (
	"regexp"
	"strings"
)

var (
	// commentRE matches semicolons at the start of a line (possibly with leading whitespace),
	// which represent comments in GCC option files.
	commentRE = regexp.MustCompile(`(?m)^[ \t]*;.*$`)

	// splitRE matches double or more newlines with any intervening whitespace,
	// separating distinct option records or metadata blocks.
	splitRE = regexp.MustCompile(`(?:\n\s*){2,}`)
)

// excludes lists metadata keywords in GCC option files that do not define command-line flags.
// See: https://gcc.gnu.org/onlinedocs/gccint/Option-file-format.html
var excludes = map[string]bool{
	"Language":       true,
	"Variable":       true,
	"TargetVariable": true,
	"TargetSave":     true,
	"HeaderKeep":     true,
	"Enum":           true,
	"EnumValue":      true,
}

// optRecord holds the raw name and properties block of an option definition.
type optRecord struct {
	name  string
	props string
}

// isOptRecord checks if the given raw record block represents an option flag,
// filtering out metadata declarations like Language, Variable, or Enum.
func isOptRecord(content string) bool {
	line, _, _ := strings.Cut(content, "\n")
	name := strings.TrimSpace(line)
	if name == "" {
		return false
	}
	return !excludes[name]
}

// extractOptRecord parses a single raw record block into its name and properties string.
func extractOptRecord(content string) optRecord {
	lines := strings.SplitN(content, "\n", 3)
	name := strings.TrimSpace(lines[0])
	var props string
	if len(lines) > 1 {
		props = strings.TrimSpace(lines[1])
	}
	return optRecord{
		name:  name,
		props: props,
	}
}

// extractOptRecords parses the raw content of a GCC .opt file, strips comments,
// and extracts all option records.
func extractOptRecords(content string) []optRecord {
	// Normalizing carriage returns prevents parsing drift when files are checked out
	// on Windows with core.autocrlf enabled.
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = commentRE.ReplaceAllString(content, "")

	contents := splitRE.Split(content, -1)
	records := make([]optRecord, 0, len(contents))
	for _, c := range contents {
		c = strings.TrimSpace(c)
		if isOptRecord(c) {
			records = append(records, extractOptRecord(c))
		}
	}
	return records
}
