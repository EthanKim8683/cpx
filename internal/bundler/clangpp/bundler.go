package clangpp

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/EthanKim8683/cpx/internal/port"
)

var bundleFlags = []string{
	"-o-",
	"-E",
	"-P",
	"-fkeep-system-includes",
	"-fdirectives-only",
}

type Bundler struct {
	executable string
	args       []string
}

func (b *Bundler) Bundle(ctx context.Context) (string, error) {
	var (
		args   = slices.Concat(b.args, bundleFlags)
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	cmd := exec.CommandContext(ctx, b.executable, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

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

func NewBundler(executable string, args []string) port.Bundler {
	return &Bundler{
		executable: executable,
		args:       args,
	}
}
