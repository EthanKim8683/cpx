package gcc

import (
	"github.com/EthanKim8683/cpx/internal/bundler/clang"
	"github.com/EthanKim8683/cpx/internal/port"
)

func NewBundler(args []string) port.Bundler {
	return clang.NewBundler(append([]string{"clang++"}, args[1:]...))
}
