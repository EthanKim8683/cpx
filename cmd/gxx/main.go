// Package main implements the g++ compiler wrapper shim.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/EthanKim8683/cpx/internal/config"
	"github.com/EthanKim8683/cpx/internal/gcc"
)

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	return gcc.RunGXX(&cfg)
}

func main() {
	if err := run(); err != nil {
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "cpx g++: %v\n", err)
		os.Exit(1)
	}
}
