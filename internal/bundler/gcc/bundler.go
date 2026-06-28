package gcc

import (
	"context"
	"path/filepath"
	"slices"
	"strings"

	"github.com/EthanKim8683/cpx/internal/bundler/clang"
	"github.com/EthanKim8683/cpx/internal/port"
)

var cppExtensions = []string{
	".cpp", ".cppm", ".hpp",
	".cc", ".hh",
	".cxx", ".cxxm", ".hxx",
	".C",
	".c++", ".c++m", ".h++",
	".ixx",
}

// detectCPP guesses whether the given command indicates a C++ compilation.
//
// Since C++ is a superset of C, the heuristic assumes that if any trace of C++
// is found in the command, it's a C++ compilation.
//
// GCC also has its own heuristics for detecting C++, but they're hard to
// implement without reimplementing GCC's option parsing.
func detectCPP(bin string, args []string) bool {
	if strings.Contains(filepath.Clean(bin), "++") {
		return true
	}
	for _, arg := range args {
		if strings.Contains(arg, "++") {
			return true
		}
		for _, ext := range cppExtensions {
			if strings.HasSuffix(arg, ext) {
				return true
			}
		}
	}
	return false
}

// Bundler wraps clang.Bundler to bundle C/C++ source files on behalf of GCC.
type Bundler struct {
	bundler port.Bundler
}

// Bundle delegates to the underlying bundler.
func (b *Bundler) Bundle(ctx context.Context) (string, error) {
	return b.bundler.Bundle(ctx)
}

var _ port.Bundler = (*Bundler)(nil)

// NewBundler adapts the given command for clang.NewBundler to emulate GCC's
// preprocessor environment.
func NewBundler(clangBin string, gccBin string, args []string) port.Bundler {
	if detectCPP(gccBin, args) {
		args = slices.Concat(gccCPPFlags, args)
	} else {
		args = slices.Concat(gccCFlags, args)
	}
	return &Bundler{
		bundler: clang.NewBundler(clangBin, args),
	}
}
