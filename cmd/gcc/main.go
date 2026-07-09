// Package main implements the gcc compiler wrapper shim.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/EthanKim8683/cpx/internal/config"
	"github.com/EthanKim8683/cpx/internal/gcc"
)

func execute() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	return gcc.ExecuteGCC(&cfg, os.Args)
}

func main() {
	if err := execute(); err != nil {
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "cpx gcc: %v\n", err)
		os.Exit(1)
	}
}
