package cdb

import (
	"fmt"
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
	Bin string
}

// Compile executes the compiler binary as a subprocess with the provided arguments.
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
