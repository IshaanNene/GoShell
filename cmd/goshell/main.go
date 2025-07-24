
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/IshaanNene/GoShell/internal/core"
	"github.com/peterh/liner"
)

func main() {
	
	historyManager, err := core.NewHistoryManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing history: %v\n", err)
	}
	defer historyManager.Close()

	
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetCompleter(func(line string) (c []string) {
		
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
		
		
		registeredCommands := core.GetRegisteredCommands()
		for _, cmd := range registeredCommands {
			if strings.HasPrefix(cmd, strings.Fields(line)[0]) {
				c = append(c, cmd)
			}
		}
		
		return
	})

	
	for _, h := range historyManager.GetHistory() {
		line.AppendHistory(h)
	}

	
	executor := core.NewExecutor()
	executor.SetHistoryManager(historyManager)

	
	
	core.InitializeCommands()

	
	missing := core.ValidateCommands()
	if len(missing) > 0 {
		fmt.Printf("Warning: Some commands not initialized: %v\n", missing)
	}

	fmt.Printf("GoShell v1.0 - Get Shelled\n")
	fmt.Printf("Type 'help' for available commands or 'exit' to quit.\n\n")

	for {
		
		input, err := line.Prompt("goshell> ")
		if err != nil {
			if err == liner.ErrPromptAborted {
				
				continue
			}
			break 
		}

		
		trimmedInput := strings.TrimSpace(input)
		if trimmedInput == "exit()" || trimmedInput == "exit" {
			break
		}

		
		if trimmedInput == "" {
			continue
		}

		
		if trimmedInput == "help" {
			showHelp()
			continue
		}

		
		line.AppendHistory(input)
		historyManager.Add(input)

		
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

	
	historyManager.Flush()
}

func showHelp() {
	fmt.Println("GoShell - Available Commands:")
	fmt.Println()
	
	
	customCommands := core.GetAvailableCommands()
	if len(customCommands) > 0 {
		fmt.Println("Custom Commands:")
		for _, cmd := range customCommands {
			fmt.Printf("  %s\n", cmd)
		}
		fmt.Println()
	}
	
	
	fmt.Println("Built-in Commands:")
	builtIns := []string{"pwd", "echo", "history", "exit", "help"}
	for _, cmd := range builtIns {
		fmt.Printf("  %s\n", cmd)
	}
	
	
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
