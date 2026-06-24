package gpp

import (
	"context"
	"errors"
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

func NewBundler(cfg config.Config, args []string) (port.Bundler, error) {
	if len(args) == 0 {
		return nil, errors.New("no arguments provided")
	}

	b, _ := clangpp.NewBundler(slices.Concat([]string{cfg.Clangpp}, flags, args[1:]))
	return &Bundler{
		clangpp: b,
	}, nil
}
