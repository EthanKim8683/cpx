package cdb

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/sync/errgroup"
)

// Shim coordinates compiling a command while concurrently recording it to
// a compilation database store.
type Shim struct {
	Name  string
	Bin   string
	Cfg   *Config
	Store *Store
}

func (s *Shim) update(args []string) error {
	command, err := Parse(s.Cfg, args)
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

	if err := s.Store.Add(records); err != nil {
		return fmt.Errorf("updating compilation database: %w", err)
	}
	return nil
}

func (s *Shim) compile(args []string) error {
	cmd := exec.Command(s.Bin, args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}
	return nil
}

// Execute runs the actual compiler binary with the provided arguments
// and concurrently parses and updates the compilation database.
func (s *Shim) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no compiler arguments provided")
	}

	var g errgroup.Group
	g.Go(func() error {
		return s.update(args)
	})

	var errs error
	errs = errors.Join(errs, s.compile(args))
	errs = errors.Join(errs, g.Wait())
	return errs
}
