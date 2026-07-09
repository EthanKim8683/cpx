package cdb

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"golang.org/x/sync/errgroup"
)

type Compiler interface {
	Compile(argv []string) error
}

type RecordAdder interface {
	Add(records []Record) error
}

// ExecCompiler implements Compiler by executing an external subprocess.
type ExecCompiler struct {
	Bin    string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (c *ExecCompiler) Compile(argv []string) error {
	cmd := exec.Command(c.Bin, argv[1:]...)
	cmd.Stdin = c.Stdin
	if cmd.Stdin == nil {
		cmd.Stdin = os.Stdin
	}
	cmd.Stdout = c.Stdout
	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = c.Stderr
	if cmd.Stderr == nil {
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}
	return nil
}

// Shim coordinates compiling a command while concurrently recording it to
// a compilation database store.
type Shim struct {
	Name        string
	Cfg         *Config
	Compiler    Compiler
	RecordAdder RecordAdder
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

	if err := s.RecordAdder.Add(records); err != nil {
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
