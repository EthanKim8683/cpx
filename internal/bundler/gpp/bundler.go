package gpp

import (
	"github.com/EthanKim8683/cpx/internal/bundler/clang"
	"github.com/EthanKim8683/cpx/internal/port"
)

func NewBundler(command string, flags []string) (port.Bundler, error) {
	return clang.NewBundler(command, flags), nil
}
