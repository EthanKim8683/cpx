package main

import (
	"io"
	"os"
)

func main() {
	_, _ = io.Copy(os.Stdout, os.Stdin)
}
