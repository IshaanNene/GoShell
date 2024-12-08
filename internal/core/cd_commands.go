package core

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

var CdCmd = &cobra.Command{
	Use:   "cd [directory]",
	Short: "Change the current working directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetDir := args[0]
		if targetDir == "~" {
			usr, err := user.Current()
			if err != nil {
				log.Fatalf("Error getting user information: %v", err)
			}
			targetDir = usr.HomeDir
		} else if targetDir == ".." {
			currentDir, err := os.Getwd()
			if err != nil {
				log.Fatalf("Error getting current directory: %v", err)
			}
			targetDir = filepath.Dir(currentDir)
		}
		err := os.Chdir(targetDir)
		if err != nil {
			log.Fatalf("Error changing directory: %v", err)
		}
		newDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting new directory: %v", err)
		}

		fmt.Println("Directory changed successfully to:", newDir)
	},
}
