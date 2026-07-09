package cdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeRecords(t *testing.T) {
	t.Parallel()

	a := []Record{
		{File: "main.cpp", Dir: "/dir1", Shim: "g++"},
		{File: "helper.cpp", Dir: "/dir1", Shim: "g++"},
	}
	b := []Record{
		{File: "main.cpp", Dir: "/dir2", Shim: "g++"},
		{File: "solve.cpp", Dir: "/dir2", Shim: "g++"},
	}

	merged := mergeRecords(a, b)
	assert.Len(t, merged, 3)

	m := make(map[string]Record)
	for _, r := range merged {
		m[r.File] = r
	}

	assert.Equal(t, "/dir2", m["main.cpp"].Dir)   // Overwritten by b
	assert.Equal(t, "/dir1", m["helper.cpp"].Dir) // Preserved from a
	assert.Equal(t, "/dir2", m["solve.cpp"].Dir)  // Added from b
}
