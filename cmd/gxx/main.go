// Package main implements the g++ compiler wrapper shim.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/EthanKim8683/cpx/internal/cdb"
	"github.com/EthanKim8683/cpx/internal/config"
	"github.com/EthanKim8683/cpx/internal/gcc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cpx g++: %v\n", err)
		os.Exit(1)
	}

	if cfg.GXX == "" {
		fmt.Fprintf(os.Stderr, "cpx g++: g++ path is not configured (set GXX environment variable)\n")
		os.Exit(1)
	}

	shim := cdb.Shim{
		Name:  "g++",
		Bin:   cfg.GXX,
		Cfg:   gcc.CDBConfig,
		Store: cdb.NewStore(cfg.CDB),
	}
	if err := shim.Execute(os.Args); err != nil {
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "cpx g++: %v\n", err)
		os.Exit(1)
	}
}
