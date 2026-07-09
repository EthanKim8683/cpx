// Package gcc provides option configuration and execution shim wrappers for the GCC toolchain.
package gcc

import (
	"errors"

	"github.com/EthanKim8683/cpx/internal/cdb"
	"github.com/EthanKim8683/cpx/internal/config"
)

const (
	gcc = "gcc"
	gxx = "g++"
)

// ExecuteGCC executes the gcc driver shim under the given configuration.
func ExecuteGCC(cfg *config.Config, args []string) error {
	bin := cfg.GCC
	if bin == "" {
		return errors.New("GCC not set")
	}

	shim := &cdb.Shim{
		Name:        gcc,
		Cfg:         CDBConfig,
		Compiler:    &cdb.ExecCompiler{Bin: bin},
		RecordAdder: cdb.NewStore(cfg.CDB),
	}
	return shim.Execute(args)
}

// ExecuteGXX executes the g++ driver shim under the given configuration.
func ExecuteGXX(cfg *config.Config, args []string) error {
	bin := cfg.GXX
	if bin == "" {
		return errors.New("GXX not set")
	}

	shim := &cdb.Shim{
		Name:        gxx,
		Cfg:         CDBConfig,
		Compiler:    &cdb.ExecCompiler{Bin: bin},
		RecordAdder: cdb.NewStore(cfg.CDB),
	}
	return shim.Execute(args)
}
