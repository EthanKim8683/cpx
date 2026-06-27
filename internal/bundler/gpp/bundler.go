package gpp

import (
	"context"

	"github.com/EthanKim8683/cpx/internal/cdb"
	"github.com/EthanKim8683/cpx/internal/port"
)

type Bundler struct {
	clangpp port.Bundler
}

func (b *Bundler) Bundle(ctx context.Context) (string, error) {
	return b.clangpp.Bundle(ctx)
}

var _ port.Bundler = (*Bundler)(nil)

func NewBundler(co cdb.CommandObject) port.Bundler {
	// return &Bundler{
	// 	clangpp: clangpp.NewBundler(executable, slices.Concat(flags, args)),
	// }
	return nil
}
