package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/EthanKim8683/cpx/internal/bundler/gcc"
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

		b := gcc.NewBundler(compileArgs)
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
