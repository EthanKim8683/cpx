//go:build ignore

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/EthanKim8683/cpx/internal/config"
	"github.com/joho/godotenv"
)

var includeRegexp = regexp.MustCompile(`(?m)^\s(\S+)$`)

func extractDefines(ctx context.Context, executable string) (map[string]string, error) {
	output, err := exec.CommandContext(ctx, executable, "-E", "-dM", "-x", "c++", os.DevNull).Output()
	if err != nil {
		return nil, fmt.Errorf("executing command: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	defines := make(map[string]string, len(lines))
	for _, line := range lines {
		tokens := strings.SplitN(line, " ", 3)
		if len(tokens) == 3 {
			defines[tokens[1]] = strings.ReplaceAll(strings.ReplaceAll(tokens[2], `\`, `\\`), `"`, `\"`)
		}
	}
	return defines, nil
}

func extractIncludes(ctx context.Context, executable string) ([]string, error) {
	output, err := exec.CommandContext(ctx, executable, "-v", "-E", "-x", "c++", os.DevNull).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("executing command: %w", err)
	}

	matches := includeRegexp.FindAllStringSubmatch(string(output), -1)
	includes := make([]string, 0, len(matches))
	for _, match := range matches {
		includes = append(includes, match[1])
	}
	return includes, nil
}

func main() {
	godotenv.Load("../../../.env")
	cfg := config.Load()
	ctx := context.Background()

	var errs error
	defines, err := extractDefines(ctx, cfg.Gpp)
	if err != nil {
		errs = errors.Join(errs, fmt.Errorf("extracting defines: %w", err))
	}
	includes, err := extractIncludes(ctx, cfg.Gpp)
	if err != nil {
		errs = errors.Join(errs, fmt.Errorf("extracting includes: %w", err))
	}
	if errs != nil {
		log.Fatal(errs.Error())
	}

	var b bytes.Buffer
	b.WriteString(`package gpp

var flags = []string{
`)
	fmt.Fprintf(&b, "\t\"-undef\",\n")
	for identifier, value := range defines {
		fmt.Fprintf(&b, "\t\"-D%s=%s\",\n", identifier, value)
	}
	for _, include := range includes {
		fmt.Fprintf(&b, "\t\"-isystem%s\",\n", include)
	}
	b.WriteString("}")
	os.WriteFile("generated_flags.go", b.Bytes(), 0644)
}
