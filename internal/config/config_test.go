package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	t.Setenv("GCC", "/usr/bin/gcc")
	t.Setenv("CLANG", "/usr/bin/clang")
	t.Setenv("CLANG_TBLGEN", "/usr/bin/clang-tblgen")

	if _, err := Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
}
