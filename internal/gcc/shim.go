// Package gcc provides option configuration and execution shim wrappers for the GCC toolchain.
package gcc

import (
	"fmt"
	"os"

	"github.com/EthanKim8683/cpx/internal/cdb"
	"github.com/EthanKim8683/cpx/internal/config"
)

const (
	gcc string = "gcc"
	gxx string = "g++"
)

// RunGCC executes the gcc driver shim under the given configuration.
func RunGCC(cfg *config.Config) error {
	bin := cfg.GCC
	if bin == "" {
		return fmt.Errorf("GCC not set")
	}

	return (&cdb.Shim{
		Name: gcc,
		Cfg:  CDBConfig,
		Compiler: &cdb.ExecCompiler{
			Bin:    bin,
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		},
		Recorder: cdb.NewFileRecorder(cfg.CDB),
	}).Execute(os.Args)
}

// RunGXX executes the g++ driver shim under the given configuration.
func RunGXX(cfg *config.Config) error {
	bin := cfg.GXX
	if bin == "" {
		return fmt.Errorf("GXX not set")
	}

	return (&cdb.Shim{
		Name: gxx,
		Cfg:  CDBConfig,
		Compiler: &cdb.ExecCompiler{
			Bin:    bin,
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		},
		Recorder: cdb.NewFileRecorder(cfg.CDB),
	}).Execute(os.Args)
}
