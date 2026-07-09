package cdb

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
)

// Shim coordinates compiling a command while concurrently recording it to
// a compilation database store.
type Shim struct {
	// Name is the name of the compiler shim (e.g., "g++").
	Name string
	// Cfg is the compiler-specific option pattern configuration.
	Cfg *Config
	// Compiler is the compiler execution dependency.
	Compiler Compiler
	// Recorder is the compilation database storage dependency.
	Recorder Recorder
}

func (s *Shim) update(argv []string) error {
	command, err := Parse(s.Cfg, argv)
	if err != nil {
		return fmt.Errorf("parsing compiler arguments: %w", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	records := make([]Record, 0, len(command.Args))
	for _, arg := range command.Args {
		records = append(records, Record{
			File:    arg,
			Dir:     dir,
			Shim:    s.Name,
			Command: command,
		})
	}

	if err := s.Recorder.Record(records); err != nil {
		return fmt.Errorf("updating compilation database: %w", err)
	}
	return nil
}

// Execute runs the actual compiler binary with the provided arguments
// and concurrently parses and updates the compilation database.
func (s *Shim) Execute(argv []string) error {
	if len(argv) == 0 {
		return fmt.Errorf("no compiler arguments provided")
	}

	var g errgroup.Group
	g.Go(func() error {
		return s.update(argv)
	})

	var errs error
	errs = errors.Join(errs, s.Compiler.Compile(argv))
	errs = errors.Join(errs, g.Wait())
	return errs
}
