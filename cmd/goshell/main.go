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
		return
	})

	
	for _, h := range historyManager.GetHistory() {
		line.AppendHistory(h)
	}

	executor := core.NewExecutor()

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