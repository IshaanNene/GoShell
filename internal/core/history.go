package core

import (
	"bufio"
	"os"
	"path/filepath"
)

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

	// Try to open existing history file
	file, err := os.OpenFile(historyPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	hm.file = file

	// Load existing history
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hm.history = append(hm.history, scanner.Text())
	}

	return hm, nil
}

func (hm *HistoryManager) Add(command string) {
	if command != "" && (len(hm.history) == 0 || hm.history[len(hm.history)-1] != command) {
		hm.history = append(hm.history, command)
	}
}

func (hm *HistoryManager) GetHistory() []string {
	return hm.history
}

func (hm *HistoryManager) Flush() error {
	if hm.file == nil {
		return nil
	}

	// Truncate file and write all history
	err := hm.file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = hm.file.Seek(0, 0)
	if err != nil {
		return err
	}

	for _, cmd := range hm.history {
		_, err = hm.file.WriteString(cmd + "
")
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
