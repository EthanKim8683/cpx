package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Setenv("GCC", "/usr/bin/gcc")
	t.Setenv("CLANG", "/usr/bin/clang")

	_, err := Load()
	require.NoError(t, err)
}
