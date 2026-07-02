package main

import (
	"regexp"
	"strings"
)

var (
	commentRE = regexp.MustCompile(`(?m)^[ \t]*;.*$`)
	splitRE   = regexp.MustCompile(`\n{2,}`)
)

// exclude records whose first line matches the excludes
// https://gcc.gnu.org/onlinedocs/gccint/Option-file-format.html
var excludes = map[string]bool{
	"Language":       true,
	"Variable":       true,
	"TargetVariable": true,
	"TargetSave":     true,
	"HeaderKeep":     true,
	"Enum":           true,
	"EnumValue":      true,
}

type optRecord struct {
	name  string
	attrs string
}

// isOptRecord determines if a raw record block defines a compiler option.
// It returns false if the record is empty or matches an excluded keyword (e.g. Language, Variable).
func isOptRecord(content string) bool {
	line, _, _ := strings.Cut(content, "\n")
	name := strings.TrimSpace(line)
	if name == "" {
		return false
	}
	return !excludes[name]
}

// extractOptRecord parses a single raw record block into an optRecord.
func extractOptRecord(content string) optRecord {
	lines := strings.SplitN(content, "\n", 3)
	name := strings.TrimSpace(lines[0])
	var attrs string
	if len(lines) > 1 {
		attrs = strings.TrimSpace(lines[1])
	}
	return optRecord{
		name:  name,
		attrs: attrs,
	}
}

// extractOptRecords parses the raw content of an option file, stripping comments
// and extracting all valid option records.
func extractOptRecords(content string) []optRecord {
	// Normalize Windows line endings
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
