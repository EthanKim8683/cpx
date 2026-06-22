package cbundler

import (
	"errors"
	"fmt"
	"path/filepath"
)

type fileNode struct {
	absPath    string
	fragment   string
	dependents []*fileNode
}

func buildGraph(absPath string, includePaths []string) ([]*fileNode, error) {
	absPath = filepath.Clean(absPath)

	var stack []string
	pushed := make(map[string]struct{})
	push := func(absPath string) {
		if _, ok := pushed[absPath]; ok {
			return
		}
		pushed[absPath] = struct{}{}
		stack = append(stack, absPath)
	}
	pop := func() string {
		absPath := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return absPath
	}

	var (
		nodeMap       = make(map[string]*fileNode)
		dependentsMap = make(map[string][]string)
	)
	var errs error
	push(absPath)
	for len(stack) > 0 {
		absPath := pop()

		fragment, dependencies, err := resolveFile(absPath, includePaths)
		errs = errors.Join(errs, err)

		nodeMap[absPath] = &fileNode{
			absPath:  absPath,
			fragment: fragment,
		}
		for _, dependency := range dependencies {
			dependentsMap[dependency] = append(dependentsMap[dependency], absPath)
		}

		for _, dependency := range dependencies {
			push(dependency)
		}
	}
	if errs != nil {
		return nil, errs
	}

	for absPath, node := range nodeMap {
		dependents := make([]*fileNode, 0, len(dependentsMap[absPath]))
		for _, dependent := range dependentsMap[absPath] {
			dependents = append(dependents, nodeMap[dependent])
		}
		node.dependents = dependents
	}

	nodes := make([]*fileNode, 0, len(nodeMap))
	for _, node := range nodeMap {
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func topologicalSort(nodes []*fileNode) ([]*fileNode, error) {
	var stack []*fileNode
	push := func(node *fileNode) {
		stack = append(stack, node)
	}
	pop := func() *fileNode {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return node
	}

	indegree := make(map[*fileNode]int)
	for _, node := range nodes {
		indegree[node] = 0
	}
	for _, node := range nodes {
		for _, dependent := range node.dependents {
			indegree[dependent]++
		}
	}

	for _, node := range nodes {
		if indegree[node] == 0 {
			push(node)
		}
	}
	var sorted []*fileNode
	for len(stack) > 0 {
		node := pop()

		sorted = append(sorted, node)

		for _, dependent := range node.dependents {
			indegree[dependent]--
			if indegree[dependent] == 0 {
				push(dependent)
			}
		}
	}
	if len(sorted) != len(nodes) {
		return nil, fmt.Errorf("cycle detected")
	}
	return sorted, nil
}
