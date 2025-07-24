package core

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var EchoCmd = &cobra.Command{
	Use:   "echo [string to echo]",
	Short: "Prints the provided string to standard output",
	Long:  `echo is a command that outputs the strings it is being passed as arguments.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(strings.Join(args, " "))
	},
}

func init() {
	// This will be updated later to use a command registry
}
