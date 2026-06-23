package gxx

import (
	"fmt"
	"os/exec"

	"github.com/EthanKim8683/cpx/internal/bundler/clangxx"
	"github.com/EthanKim8683/cpx/internal/port"
)

func NewBundler(flags []string) (port.Bundler, error) {
	if _, err := exec.LookPath("clang++"); err != nil {
		return nil, fmt.Errorf("missing clang++: %w", err)
	}
	return clangxx.NewBundler("clang++", flags), nil
}
