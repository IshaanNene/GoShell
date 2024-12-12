package main

import (
	"os"

	"github.com/IshaanNene/GoShell/internal/core"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "goshell",
		Short: "A simple shell command executor",
		Run: func(cmd *cobra.Command, args []string) {
			color.Green("Welcome to GoShell! Use 'goshell help' to see available commands.")
		},
	}
	rootCmd.AddCommand(core.LsCmd)
	rootCmd.AddCommand(core.CdCmd)
	rootCmd.AddCommand(core.MkdirCmd)
	rootCmd.AddCommand(core.RmCmd)
	rootCmd.AddCommand(core.TouchCmd)
	rootCmd.AddCommand(core.PwdCmd)
	rootCmd.AddCommand(core.CatCmd)
	rootCmd.AddCommand(core.DateCmd)
	rootCmd.AddCommand(core.Iamwho)

	if err := rootCmd.Execute(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}
