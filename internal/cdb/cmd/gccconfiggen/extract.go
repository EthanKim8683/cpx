// extract.go contains logic for extracting option records from GCC option files.

package main

import (
	"regexp"
	"strings"
)

var (
	commentRE = regexp.MustCompile(`(?m)^[ \t]*;.*$`)
	splitRE   = regexp.MustCompile(`(?:\n\s*){2,}`)
)

// excludes contains record types to exclude when extracting option records.
// Record types documentation:
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

// optRecord represents an extracted option record.
type optRecord struct {
	name  string
	props string
}

// isOptRecord checks if the given record block defines an option record.
func isOptRecord(content string) bool {
	line, _, _ := strings.Cut(content, "\n")
	name := strings.TrimSpace(line)
	if name == "" {
		return false
	}
	return !excludes[name]
}

// extractOptRecord parses a record block into an optRecord.
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

// extractOptRecords extracts option records from the content of a GCC .opt file.
func extractOptRecords(content string) []optRecord {
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
