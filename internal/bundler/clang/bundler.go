package clang

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
	"-E",                     // Run preprocessor only
	"-P",                     // Disable linemarkers in output
	"-fkeep-system-includes", // Only expand user includes
	"-fdirectives-only",      // Keep macros intact
}

// removeOutputFlags removes output flags (e.g. -ofoo, -o foo, --output=foo,
// --output foo) from the argument list so the preprocessor outputs to stdout
// instead.
//
// Appending -o- to the argument list (overriding all prior output flags)
// achieves a similar result, but is less explicit and causes Clang to complain
// about -o being set when bundling multiple files.
func removeOutputFlags(args []string) []string {
	filtered := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-o":
			i++
		case arg == "--output":
			i++
		// Some flags (e.g. -openmp, -output-pch) and filenames may start with -o,
		// but these are rare in practice.
		case strings.HasPrefix(arg, "-o"):
		case strings.HasPrefix(arg, "--output="):
		default:
			filtered = append(filtered, arg)
		}
	}
	return filtered
}

// Bundler uses the Clang preprocessor to bundle C/C++ source files.
type Bundler struct {
	bin  string
	args []string
}

// Bundle runs the bundle command and returns stdout, or an error containing
// stderr if it fails.
func (b *Bundler) Bundle(ctx context.Context) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, b.bin, b.args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		err = fmt.Errorf("bundling with clang++ preprocessor: %w", err)
		if reason := strings.TrimSpace(stderr.String()); reason != "" {
			err = fmt.Errorf("%w: %s", err, reason)
		}
		return "", err
	}
	return stdout.String(), nil
}

var _ port.Bundler = (*Bundler)(nil)

// NewBundler transforms the given command into a bundle command and creates a
// Bundler to run it.
func NewBundler(bin string, args []string) port.Bundler {
	// The transformations below are not exhaustive; flags like -c, -v, -dM, etc.
	// interfere with bundling, but we can safely assume they're out of the scope
	// of competitive programming.
	args = removeOutputFlags(args)
	args = slices.Concat(args, bundleFlags)
	return &Bundler{
		bin:  bin,
		args: args,
	}
}
