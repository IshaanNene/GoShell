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
	if NetCmd != nil {
		RegisterCommand(NetCmd)
	}
	if MemCmd != nil {
		RegisterCommand(MemCmd)
	}
	if DiskCmd != nil {
		RegisterCommand(DiskCmd)
	}
	if CPUCmd != nil {
		RegisterCommand(CPUCmd)
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
	if NetCmd == nil {
		missing = append(missing, "net")
	}
	if MemCmd == nil {
		missing = append(missing, "mem")
	}
	if DiskCmd == nil {
		missing = append(missing, "disk")
	}
	if CPUCmd == nil {
		missing = append(missing, "cpu")
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
	if NetCmd != nil {
		commands = append(commands, "net")
	}
	if MemCmd != nil {
		commands = append(commands, "mem")
	}
	if DiskCmd != nil {
		commands = append(commands, "disk")
	}
	if CPUCmd != nil {
		commands = append(commands, "cpu")
	}

	return commands
}

func InitCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(LsCmd)
	rootCmd.AddCommand(CatCmd)
	rootCmd.AddCommand(RmCmd)
	rootCmd.AddCommand(CdCmd)
	rootCmd.AddCommand(TouchCmd)
	rootCmd.AddCommand(MkdirCmd)
	rootCmd.AddCommand(NetCmd)
	rootCmd.AddCommand(MemCmd)
	rootCmd.AddCommand(DiskCmd)
	rootCmd.AddCommand(CPUCmd)
}
