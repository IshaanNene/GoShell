package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/IshaanNene/GoShell/internal/core"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
)

func main() {
	color.Green("Welcome to GoShell! Press Ctrl+C to exit.")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ") // Simple prompt

		input, _ := reader.ReadString('
')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		// Here we will add the parsing and execution logic
		// For now, we just print the input
		fmt.Printf("You entered: %s
", input)

		// In a real shell, you would parse the input,
		// execute the command, and handle output.
	}
}

// We will keep the rootCmd and command definitions for now,
// but they will be used differently later when we implement
// the command parsing and execution.
var rootCmd = &cobra.Command{
	Use:   "goshell",
	Short: "A simple shell command executor",
	Run: func(cmd *cobra.Command, args []string) {
		color.Green("Welcome to GoShell! Use 'goshell help' to see available commands.")
	},
}

func init() {
	rootCmd.AddCommand(core.LsCmd)
	rootCmd.AddCommand(core.CdCmd)
	rootCmd.AddCommand(core.MkdirCmd)
	rootCmd.AddCommand(core.RmCmd)
	rootCmd.AddCommand(core.TouchCmd)
	rootCmd.AddCommand(core.PwdCmd)
	rootCmd.AddCommand(core.CatCmd)
	rootCmd.AddCommand(core.DateCmd)
	rootCmd.AddCommand(core.Iamwho)
}
