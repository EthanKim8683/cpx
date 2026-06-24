package gcc

import (
	"context"

	"github.com/EthanKim8683/cpx/internal/bundler/clang"
	"github.com/EthanKim8683/cpx/internal/port"
)

type Bundler struct {
	clangBundler port.Bundler
}

func (b *Bundler) Bundle(ctx context.Context) (string, error) {
	return b.clangBundler.Bundle(ctx)
}

func NewBundler(args []string) port.Bundler {
	clangArgs := []string{"clang++"}
	clangArgs = append(clangArgs, clangFlags...)
	clangArgs = append(clangArgs, args[1:]...)
	return &Bundler{
		clangBundler: clang.NewBundler(clangArgs),
	}
}
