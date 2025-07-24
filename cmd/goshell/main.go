// main.go
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
		fmt.Fprintf(os.Stderr, "Error initializing history: %v
", err)
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
			// You can add more sophisticated completion logic here
			files, _ := os.ReadDir("./")
			for _, f := range files {
				if strings.HasPrefix(f.Name(), prefix) {
					c = append(c, f.Name())
				}
			}
		}
		return
	})

	// Load history into liner
	for _, h := range historyManager.GetHistory() {
		line.AppendHistory(h)
	}

	executor := core.NewExecutor()

	for {
		// Get current working directory for the prompt
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v
", err)
			wd = "?"
		}
		prompt := fmt.Sprintf("%s> ", wd)

		// Read input
		input, err := line.Prompt(prompt)
		if err != nil {
			if err == liner.ErrPromptAborted {
				// User pressed Ctrl+C
				continue
			}
			break // Exit on other errors (e.g., EOF)
		}

		// Add to history
		line.AppendHistory(input)
		historyManager.Add(input)

		// Parse and execute
		if input != "" {
			chain, err := core.ParseInput(input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v
", err)
				continue
			}

			if chain != nil {
				if err := executor.ExecuteChain(chain); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v
", err)
				}
			}
		}
	}

	// Save history
	historyManager.Flush()
}
