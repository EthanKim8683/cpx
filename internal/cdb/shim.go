package cdb

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/sync/errgroup"
)

// Compiler runs a compile command.
//
// ExecCompiler is the real implementation. Shim depends on Compiler so tests
// can substitute a fake.
type Compiler interface {
	Compile(argv []string) error
}

// ExecCompiler is the subprocess-backed Compiler implementation.
type ExecCompiler struct {
	// Bin is the absolute path to the compiler driver executable.
	Bin string
}

// Compile runs Bin with argv[1:]. Stdin, stdout, and stderr are inherited from
// the shim process so behavior matches invoking the driver directly.
func (c *ExecCompiler) Compile(argv []string) error {
	//nolint:gosec // c.Bin and argv are external compiler driver inputs we must execute
	cmd := exec.Command(c.Bin, argv[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}
	return nil
}

var _ Compiler = (*ExecCompiler)(nil)

// Shim coordinates compiling a command while concurrently recording it to
// a compilation database store.
type Shim struct {
	// Name is the name of the compiler shim (e.g., "g++").
	Name string
	// Cfg is the compiler-specific option pattern configuration.
	Cfg *Config
	// Compiler is the compiler execution dependency.
	Compiler Compiler
	// Store is the compilation database storage dependency.
	Store Store
}

func (s *Shim) update(argv []string, dir string) error {
	command, err := Parse(s.Cfg, argv)
	if err != nil {
		return fmt.Errorf("parsing compiler arguments: %w", err)
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

// Execute runs the actual compiler binary with the provided arguments
// and concurrently parses and updates the compilation database.
func (s *Shim) Execute(argv []string, dir string) error {
	if len(argv) == 0 {
		return errors.New("no compiler arguments provided")
	}

	var g errgroup.Group
	g.Go(func() error {
		return s.update(argv, dir)
	})

	var errs error
	errs = errors.Join(errs, s.Compiler.Compile(argv))
	errs = errors.Join(errs, g.Wait())
	return errs
}
