// cmd/goshell/main.go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/IshaanNene/GoShell/internal/core"
	"github.com/peterh/liner"
)

func main() {
	// Initialize history
	historyManager, err := core.NewHistoryManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing history: %v\n", err)
	}
	defer historyManager.Close()

	// Initialize liner for line editing and history
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetCompleter(func(line string) (c []string) {
		// Basic file path completion
		if strings.Contains(line, " ") {
			parts := strings.Split(line, " ")
			prefix := parts[len(parts)-1]
			files, _ := os.ReadDir("./")
			for _, f := range files {
				if strings.HasPrefix(f.Name(), prefix) {
					c = append(c, f.Name())
				}
			}
		}
		
		// Add registered command completion
		registeredCommands := core.GetRegisteredCommands()
		for _, cmd := range registeredCommands {
			if strings.HasPrefix(cmd, strings.Fields(line)[0]) {
				c = append(c, cmd)
			}
		}
		
		return
	})

	// Load history into liner
	for _, h := range historyManager.GetHistory() {
		line.AppendHistory(h)
	}

	// Initialize executor and set history manager
	executor := core.NewExecutor()
	executor.SetHistoryManager(historyManager)

	// Initialize all your custom commands
	// This registers the commands that are already defined in your command files
	core.InitializeCommands()

	// Validate that commands were properly initialized
	missing := core.ValidateCommands()
	if len(missing) > 0 {
		fmt.Printf("Warning: Some commands not initialized: %v\n", missing)
	}

	fmt.Printf("GoShell v1.0 - Enhanced Shell\n")
	fmt.Printf("Registered commands: %d\n", core.GetCommandCount())
	fmt.Printf("Type 'help' for available commands or 'exit' to quit.\n\n")

	for {
		// Read input with goshell> prompt
		input, err := line.Prompt("goshell> ")
		if err != nil {
			if err == liner.ErrPromptAborted {
				// User pressed Ctrl+C
				continue
			}
			break // Exit on other errors (e.g., EOF)
		}

		// Check for exit command
		trimmedInput := strings.TrimSpace(input)
		if trimmedInput == "exit()" || trimmedInput == "exit" {
			break
		}

		// Skip empty input
		if trimmedInput == "" {
			continue
		}

		// Handle help command
		if trimmedInput == "help" {
			showHelp()
			continue
		}

		// Add to history
		line.AppendHistory(input)
		historyManager.Add(input)

		// Parse and execute
		chain, err := core.ParseInput(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}

		if chain != nil {
			if err := executor.ExecuteChain(chain); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		}
	}

	// Save history
	historyManager.Flush()
}

func showHelp() {
	fmt.Println("GoShell - Available Commands:")
	fmt.Println()
	
	// Show custom commands that are available
	customCommands := core.GetAvailableCommands()
	if len(customCommands) > 0 {
		fmt.Println("Custom Commands:")
		for _, cmd := range customCommands {
			fmt.Printf("  %s\n", cmd)
		}
		fmt.Println()
	}
	
	// Show built-in commands
	fmt.Println("Built-in Commands:")
	builtIns := []string{"pwd", "echo", "history", "exit", "help"}
	for _, cmd := range builtIns {
		fmt.Printf("  %s\n", cmd)
	}
	
	// Show registered commands from registry
	registered := core.GetRegisteredCommands()
	if len(registered) > 0 {
		fmt.Println("\nRegistered Commands:")
		for _, cmd := range registered {
			fmt.Printf("  %s\n", cmd)
		}
	}
	
	fmt.Println("\nOperators:")
	fmt.Println("  |   - Pipe output to next command")
	fmt.Println("  &&  - Execute next command if current succeeds")
	fmt.Println("  ||  - Execute next command if current fails")
	fmt.Println("  ;   - Execute commands sequentially")
	fmt.Println("  >   - Redirect output to file")
	fmt.Println("  >>  - Append output to file")
	fmt.Println("  <   - Redirect input from file")
	fmt.Println("  2>  - Redirect errors to file")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ls -la > output.txt")
	fmt.Println("  rm -i file.txt")
	fmt.Println("  ls | grep .go")
	fmt.Println("  mkdir test && cd test")
	fmt.Println()
}
