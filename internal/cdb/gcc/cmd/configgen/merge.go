package main

import (
	"slices"
	"strings"
)

// mergeRecords merges optRecords with the same name into a single record by concatenating their props.
func mergeRecords(records []optRecord) []optRecord {
	slices.SortFunc(records, func(a, b optRecord) int {
		return strings.Compare(a.name, b.name)
	})

	merged := make([]optRecord, 0, len(records))
	for i := 0; i < len(records); {
		name := records[i].name
		var b strings.Builder
		b.WriteString(records[i].props)
		for i++; i < len(records); i++ {
			if records[i].name != name {
				break
			}
			b.WriteString(" ")
			b.WriteString(records[i].props)
		}
		merged = append(merged, optRecord{
			name:  name,
			props: b.String(),
		})
	}
	return merged
}
