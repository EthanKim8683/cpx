package cdb

import (
	"fmt"
	"io"
	"os"
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
