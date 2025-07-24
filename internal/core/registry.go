
package core

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var commandRegistry = make(map[string]*cobra.Command)


func RegisterCommand(cmd *cobra.Command) {
	commandRegistry[cmd.Name()] = cmd
}


func GetCommand(name string) *cobra.Command {
	return commandRegistry[name]
}


func GetRegisteredCommands() []string {
	var commands []string
	for name := range commandRegistry {
		commands = append(commands, name)
	}
	return commands
}


func HistoryFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		
		return ".goshell_history"
	}
	return filepath.Join(homeDir, ".goshell_history")
}


func IsCommandRegistered(name string) bool {
	_, exists := commandRegistry[name]
	return exists
}


func UnregisterCommand(name string) {
	delete(commandRegistry, name)
}


func ClearRegistry() {
	commandRegistry = make(map[string]*cobra.Command)
}


func GetCommandCount() int {
	return len(commandRegistry)
}