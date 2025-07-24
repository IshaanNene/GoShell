package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const maxHistorySize = 1000

type HistoryManager struct {
	history []string
	file    *os.File
}

func NewHistoryManager() (*HistoryManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	historyPath := filepath.Join(homeDir, ".goshell_history")
	
	hm := &HistoryManager{
		history: make([]string, 0),
	}

	
	file, err := os.OpenFile(historyPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	hm.file = file

	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			hm.history = append(hm.history, line)
		}
	}

	return hm, nil
}

func (hm *HistoryManager) Add(command string) {
	command = strings.TrimSpace(command)
	if command == "" {
		return
	}
	
	
	if len(hm.history) > 0 && hm.history[len(hm.history)-1] == command {
		return
	}
	
	hm.history = append(hm.history, command)
	
	
	if len(hm.history) > maxHistorySize {
		hm.history = hm.history[len(hm.history)-maxHistorySize:]
	}
}

func (hm *HistoryManager) GetHistory() []string {
	return hm.history
}

func (hm *HistoryManager) GetHistoryItem(index int) (string, error) {
	if index < 0 || index >= len(hm.history) {
		return "", fmt.Errorf("history index out of range")
	}
	return hm.history[index], nil
}

func (hm *HistoryManager) SearchHistory(pattern string) []string {
	var matches []string
	for _, cmd := range hm.history {
		if strings.Contains(cmd, pattern) {
			matches = append(matches, cmd)
		}
	}
	return matches
}

func (hm *HistoryManager) PrintHistory(args []string) {
	start := 0
	count := len(hm.history)
	
	if len(args) > 0 {
		if n, err := strconv.Atoi(args[0]); err == nil {
			if n > 0 && n < count {
				start = count - n
			}
		}
	}
	
	for i := start; i < len(hm.history); i++ {
		fmt.Printf("%4d  %s\n", i+1, hm.history[i])
	}
}

func (hm *HistoryManager) Flush() error {
	if hm.file == nil {
		return nil
	}

	
	err := hm.file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = hm.file.Seek(0, 0)
	if err != nil {
		return err
	}

	for _, cmd := range hm.history {
		_, err = hm.file.WriteString(cmd + "\n")
		if err != nil {
			return err
		}
	}

	return hm.file.Sync()
}

func (hm *HistoryManager) Close() error {
	if hm.file != nil {
		hm.Flush()
		return hm.file.Close()
	}
	return nil
}