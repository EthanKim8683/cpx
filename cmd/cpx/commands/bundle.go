package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/EthanKim8683/cpx/internal/bundler/gpp"
	"github.com/spf13/cobra"
)

type CompileCommand struct {
	File      string   `json:"file"`
	Arguments []string `json:"arguments"`
	Directory string   `json:"directory"`
	Output    string   `json:"output"`
}

var BundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Bundle a source file",
	Run: func(cmd *cobra.Command, args []string) {
		bytes, err := os.ReadFile("./.cpx/compile_commands.json")
		if err != nil {
			log.Fatalf("reading compile_commands.json: %v", err)
		}

		var compileCommands []CompileCommand
		if err := json.Unmarshal(bytes, &compileCommands); err != nil {
			log.Fatalf("unmarshalling compile_commands.json: %v", err)
		}

		var compileArgs []string
		for _, compileCommand := range compileCommands {
			if compileCommand.File == args[0] {
				compileArgs = compileCommand.Arguments
			}
		}
		if compileArgs == nil {
			log.Fatalf("compile command not found for %s", args[0])
		}

		// TODO: kinda hacky; make this an option or something
		withoutDefines := make([]string, 0, len(compileArgs))
		for i := 0; i < len(compileArgs); i++ {
			arg := compileArgs[i]
			switch {
			case arg == "-D":
				i++
			case strings.HasPrefix(arg, "-D"):
			default:
				withoutDefines = append(withoutDefines, arg)
			}
		}

		b := gpp.NewBundler(withoutDefines)
		bundle, err := b.Bundle(context.Background())
		if err != nil {
			log.Fatalf("bundling: %v", err)
		}
		fmt.Println(bundle)
	},
}

func init() {
	RootCmd.AddCommand(BundleCmd)
}
