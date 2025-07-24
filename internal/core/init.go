package core

import (
	"github.com/spf13/cobra"
)

func InitializeCommands() {
	if LsCmd != nil {
		RegisterCommand(LsCmd)
	}
	
	if RmCmd != nil {
		RegisterCommand(RmCmd)
	}
	
	if CatCmd != nil {
		RegisterCommand(CatCmd)
	}
	
	if MkdirCmd != nil {
		RegisterCommand(MkdirCmd)
	}
	
	if TouchCmd != nil {
		RegisterCommand(TouchCmd)
	}
	
	if CdCmd != nil {
		RegisterCommand(CdCmd)
	}
	
	
	registerAdditionalCommands()
}



func registerAdditionalCommands() {
	
	
}


func ValidateCommands() []string {
	var missing []string
	
	if LsCmd == nil {
		missing = append(missing, "ls") 
	}
	if RmCmd == nil {
		missing = append(missing, "rm")
	}
	if CatCmd == nil {
		missing = append(missing, "cat")
	}
	if MkdirCmd == nil {
		missing = append(missing, "mkdir")
	}
	if TouchCmd == nil {
		missing = append(missing, "touch")
	}
	if CdCmd == nil {
		missing = append(missing, "cd")
	}
	
	return missing
}


func GetAvailableCommands() []string {
	var commands []string
	
	if LsCmd != nil {
		commands = append(commands, "ls")
	}
	if RmCmd != nil {
		commands = append(commands, "rm")
	}
	if CatCmd != nil {
		commands = append(commands, "cat")
	}
	if MkdirCmd != nil {
		commands = append(commands, "mkdir")
	}
	if TouchCmd != nil {
		commands = append(commands, "touch")
	}
	if CdCmd != nil {
		commands = append(commands, "cd")
	}
	
	return commands
}
func InitCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(LsCmd)
	rootCmd.AddCommand(CatCmd)
	rootCmd.AddCommand(RmCmd)
	rootCmd.AddCommand(CdCmd)
	rootCmd.AddCommand(DateCmd)
	rootCmd.AddCommand(PwdCmd)
	rootCmd.AddCommand(TouchCmd)
	rootCmd.AddCommand(MkdirCmd)
}
