package cbundler

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/EthanKim8683/cpx/internal/port"
)

type Bundler struct {
	includePaths []string
}

func (b *Bundler) Bundle(sourcePath string) (string, error) {
	if !filepath.IsAbs(sourcePath) {
		return "", fmt.Errorf("source path is not absolute: %s", sourcePath)
	}

	nodes, err := buildGraph(sourcePath, b.includePaths)
	if err != nil {
		return "", err
	}

	sortedNodes, err := topologicalSort(nodes)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for i, node := range sortedNodes {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(node.fragment)
	}
	return sb.String(), nil
}

var _ port.Bundler = (*Bundler)(nil)

func New(includePaths []string) port.Bundler {
	return &Bundler{
		includePaths: includePaths,
	}
}
