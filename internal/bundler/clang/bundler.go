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
	args []string
}

func (b *Bundler) Bundle(ctx context.Context) (string, error) {
	var (
		command = b.args[0]
		args    = append(b.args[1:],
			"-o-",
			"-E",
			"-P",
			"-fkeep-system-includes",
			"-fdirectives-only",
		)
		stdout = bytes.NewBuffer([]byte{})
		stderr = bytes.NewBuffer([]byte{})
	)
	cmd := exec.CommandContext(ctx, command, args...)
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

func NewBundler(args []string) *Bundler {
	return &Bundler{
		args: args,
	}
}
