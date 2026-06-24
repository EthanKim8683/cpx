package gpp

import (
	"context"

	"github.com/EthanKim8683/cpx/internal/bundler/clangpp"
	"github.com/EthanKim8683/cpx/internal/config"
	"github.com/EthanKim8683/cpx/internal/port"
)

type Bundler struct {
	bundler port.Bundler
}

func (b *Bundler) Bundle(ctx context.Context) (string, error) {
	return b.bundler.Bundle(ctx)
}

func NewBundler(cfg config.Config, args []string) port.Bundler {
	return &Bundler{
		bundler: clangpp.NewBundler(append(append([]string{cfg.Clangpp}, flags...), args[1:]...)),
	}
}
