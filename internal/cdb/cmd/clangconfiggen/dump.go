// Clang TableGen options JSON dumping implementation.

package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
)

// dumpJSON runs clang-tblgen to generate a JSON dump of the option definitions.
func dumpJSON(clangTblgenPath, dir string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(
		clangTblgenPath,
		"-I", filepath.Join(dir, "llvm/include"),
		"-I", filepath.Join(dir, "clang/include/clang/Options"),
		"-dump-json",
		filepath.Join(dir, "clang/include/clang/Options/Options.td"),
		"-o", "-",
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("clang-tblgen failed: %w\nstderr:\n%s", err, stderr.String())
	}
	return stdout.Bytes(), nil
}
