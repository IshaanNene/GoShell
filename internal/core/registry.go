package core

import "github.com/spf13/cobra"

var commandRegistry = make(map[string]*cobra.Command)

// RegisterCommand registers a new command.
func RegisterCommand(cmd *cobra.Command) {
	commandRegistry[cmd.Name()] = cmd
}

// GetCommand returns a command from the registry.
func GetCommand(name string) *cobra.Command {
	return commandRegistry[name]
}

// GetRegisteredCommands returns a list of registered command names.
func GetRegisteredCommands() []string {
	var commands []string
	for name := range commandRegistry {
		commands = append(commands, name)
	}
	return commands
}

func HistoryFilePath() string {
	return historyFilePath
}
