package cbundler

import (
	"errors"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func graphTestdataRoot(t *testing.T) string {
	t.Helper()

	return filepath.Join(testdataRoot(t), "graph")
}

func cyclicNodesFixture(t *testing.T) []*fileNode {
	t.Helper()

	r := filepath.Join(graphTestdataRoot(t), "cyclic")

	foo := &fileNode{
		absPath:  filepath.Join(r, "foo.cpp"),
		fragment: `// #include "bar.hpp"`,
	}
	bar := &fileNode{
		absPath:  filepath.Join(r, "bar.hpp"),
		fragment: `// #include "baz.hpp"`,
	}
	baz := &fileNode{
		absPath:  filepath.Join(r, "baz.hpp"),
		fragment: `// #include "qux.hpp"`,
	}
	qux := &fileNode{
		absPath:  filepath.Join(r, "qux.hpp"),
		fragment: `// #include "bar.hpp"`,
	}
	foo.dependents = []*fileNode{}
	bar.dependents = []*fileNode{foo, qux}
	baz.dependents = []*fileNode{bar}
	qux.dependents = []*fileNode{baz}
	return []*fileNode{foo, bar, baz, qux}
}

func diamondNodesFixture(t *testing.T) []*fileNode {
	t.Helper()

	r := filepath.Join(graphTestdataRoot(t), "diamond")

	foo := &fileNode{
		absPath: filepath.Join(r, "foo.cpp"),
		fragment: `// #include "bar.hpp"
// #include "baz.hpp"`,
	}
	bar := &fileNode{
		absPath:  filepath.Join(r, "bar.hpp"),
		fragment: `// #include "qux.hpp"`,
	}
	baz := &fileNode{
		absPath:  filepath.Join(r, "baz.hpp"),
		fragment: `// #include "qux.hpp"`,
	}
	qux := &fileNode{
		absPath:  filepath.Join(r, "qux.hpp"),
		fragment: ``,
	}
	foo.dependents = []*fileNode{}
	bar.dependents = []*fileNode{foo}
	baz.dependents = []*fileNode{foo}
	qux.dependents = []*fileNode{bar, baz}
	return []*fileNode{foo, bar, baz, qux}
}

func treeNodesFixture(t *testing.T) []*fileNode {
	t.Helper()

	r := filepath.Join(graphTestdataRoot(t), "tree")

	foo := &fileNode{
		absPath:  filepath.Join(r, "foo.cpp"),
		fragment: `// #include "bar.hpp"`,
	}
	bar := &fileNode{
		absPath: filepath.Join(r, "bar.hpp"),
		fragment: `// #include "baz.hpp"
// #include "qux.hpp"`,
	}
	baz := &fileNode{
		absPath:  filepath.Join(r, "baz.hpp"),
		fragment: ``,
	}
	qux := &fileNode{
		absPath:  filepath.Join(r, "qux.hpp"),
		fragment: ``,
	}
	foo.dependents = []*fileNode{}
	bar.dependents = []*fileNode{foo}
	baz.dependents = []*fileNode{bar}
	qux.dependents = []*fileNode{bar}
	return []*fileNode{foo, bar, baz, qux}
}

func TestBuildGraph(t *testing.T) {
	t.Parallel()

	r := graphTestdataRoot(t)
	includePaths := []string{
		filepath.Join(r, "include"),
		filepath.Join(r, "include", "include"),
	}
	tests := map[string]struct {
		sourcePath string
		nodes      []*fileNode
		err        error
	}{
		"tree": {
			sourcePath: filepath.Join(r, "tree", "foo.cpp"),
			nodes:      treeNodesFixture(t),
		},
		"diamond": {
			sourcePath: filepath.Join(r, "diamond", "foo.cpp"),
			nodes:      diamondNodesFixture(t),
		},
		"cyclic": {
			sourcePath: filepath.Join(r, "cyclic", "foo.cpp"),
			nodes:      cyclicNodesFixture(t),
		},
		"broken": {
			sourcePath: filepath.Join(r, "broken", "foo.cpp"),
			err: errors.Join(
				errors.New("could not resolve include: baz.hpp"),
				errors.New("could not resolve include: qux.hpp"),
			),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			nodes, err := buildGraph(test.sourcePath, includePaths)
			if test.err != nil {
				require.EqualError(t, err, test.err.Error())
				assert.Empty(t, nodes)
			} else {
				require.NoError(t, err)

				for _, node := range nodes {
					slices.SortFunc(node.dependents, func(a, b *fileNode) int {
						return strings.Compare(a.absPath, b.absPath)
					})
				}
				assert.ElementsMatch(t, test.nodes, nodes)
			}
		})
	}
}

func TestTopologicalSort(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		nodes []*fileNode
		err   error
	}{
		"tree": {
			nodes: treeNodesFixture(t),
		},
		"diamond": {
			nodes: diamondNodesFixture(t),
		},
		"cyclic": {
			nodes: cyclicNodesFixture(t),
			err:   errors.New("cycle detected"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			sorted, err := topologicalSort(test.nodes)
			if test.err != nil {
				require.EqualError(t, err, test.err.Error())
				assert.Empty(t, sorted)
			} else {
				require.NoError(t, err)

				seen := make(map[*fileNode]struct{})
				for _, node := range sorted {
					for _, dependent := range node.dependents {
						_, ok := seen[dependent]
						assert.False(t, ok)
					}
					seen[node] = struct{}{}
				}
			}
		})
	}
}
