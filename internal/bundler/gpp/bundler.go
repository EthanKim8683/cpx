package gpp

import (
	"context"
	"slices"

	"github.com/EthanKim8683/cpx/internal/bundler/clangpp"
	"github.com/EthanKim8683/cpx/internal/config"
	"github.com/EthanKim8683/cpx/internal/port"
)

type Bundler struct {
	clangpp port.Bundler
}

func (b *Bundler) Bundle(ctx context.Context) (string, error) {
	return b.clangpp.Bundle(ctx)
}

var _ port.Bundler = (*Bundler)(nil)

func NewBundler(cfg config.Config, args []string) port.Bundler {
	return &Bundler{
		clangpp: clangpp.NewBundler(cfg.Clangpp, slices.Concat(flags, args)),
	}
}
