package clang

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/EthanKim8683/cpx/internal/port"
)

type Bundler struct {
	command string
	flags   []string
}

func (b *Bundler) Bundle(ctx context.Context, sourcePath string) (string, error) {
	var (
		stdout = bytes.NewBuffer([]byte{})
		stderr = bytes.NewBuffer([]byte{})
	)
	cmd := exec.CommandContext(ctx, b.command, append(b.flags, sourcePath)...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		err = fmt.Errorf("bundling with clang preprocessor: %w", err)
		if reason := strings.TrimSpace(stderr.String()); reason != "" {
			err = fmt.Errorf("%w: %s", err, reason)
		}
		return "", err
	}
	return stdout.String(), nil
}

var _ port.Bundler = (*Bundler)(nil)

func NewBundler(command string, flags []string) *Bundler {
	return &Bundler{
		command: command,
		flags: append(
			flags,
			"-o-",
			"-E",
			"-P",
			"-fkeep-system-includes",
			"-fdirectives-only",
		),
	}
}
