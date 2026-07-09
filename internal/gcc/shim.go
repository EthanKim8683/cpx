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

func RunGCC(cfg *config.Config) error {
	bin := cfg.GCC
	if bin == "" {
		return fmt.Errorf("GCC not set")
	}

	return (&cdb.Shim{
		Name:     gcc,
		Cfg:      CDBConfig,
		Compiler: &cdb.ExecCompiler{
			Bin:    bin,
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		},
		RecordAdder: cdb.NewFileStore(cfg.CDB),
	}).Execute(os.Args)
}

func RunGXX(cfg *config.Config) error {
	bin := cfg.GXX
	if bin == "" {
		return fmt.Errorf("GXX not set")
	}

	return (&cdb.Shim{
		Name:     gxx,
		Cfg:      CDBConfig,
		Compiler: &cdb.ExecCompiler{
			Bin:    bin,
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		},
		RecordAdder: cdb.NewFileStore(cfg.CDB),
	}).Execute(os.Args)
}
