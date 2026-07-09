package cdb

import (
	"fmt"
	"io"
	"os/exec"
)

// Compiler defines the interface for executing a compile command.
type Compiler interface {
	Compile(argv []string) error
}

// ExecCompiler implements Compiler by executing an external subprocess.
type ExecCompiler struct {
	// Bin is the path to the compiler executable.
	Bin    string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Compile executes the compiler binary as a subprocess with the provided arguments.
func (c *ExecCompiler) Compile(argv []string) error {
	cmd := exec.Command(c.Bin, argv[1:]...)
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}
	return nil
}
