package gxx

import (
	"fmt"
	"os/exec"

	"github.com/EthanKim8683/cpx/internal/bundler/clangxx"
)

func NewBundler(flags []string) (*clangxx.Bundler, error) {
	if _, err := exec.LookPath("clang++"); err != nil {
		return nil, fmt.Errorf("missing clang++: %w", err)
	}

	return clangxx.NewBundler("clang++", flags), nil
}
