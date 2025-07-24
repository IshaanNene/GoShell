package unit

import (
	"github.com/IshaanNene/GoShell/internal/core"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"testing"
)

func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	cmd.SetArgs(args)
	output, err := cmd.ExecuteC()
	if err != nil {
		return "", err
	}
	
	return output.UsageString(), nil 
}

func TestTouchCommand(t *testing.T) {
	fileName := "testfile.txt"

	if _, err := executeCommand(core.TouchCmd, fileName); err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Fatalf("Expected file %s to exist but it does not", fileName)
	}

	os.Remove(fileName)
}

func TestLsCommand(t *testing.T) {
	file1 := "file1.txt"
	file2 := "file2.txt"
	os.Create(file1)
	os.Create(file2)
	defer os.Remove(file1)
	defer os.Remove(file2)

	expected := "file1.txt\nfile2.txt\n"

	got, err := executeCommand(core.LsCmd)
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if got != expected {
		t.Errorf("Expected %q but got %q", expected, got)
	}
}

func TestPwdCommand(t *testing.T) {
	expected, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting the current directory: %v", err)
	}

	got, err := executeCommand(core.PwdCmd)
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if got != expected {
		t.Errorf("Expected %q but got %q", expected, got)
	}
}

func TestCdCommand(t *testing.T) {
	initialDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting the initial directory: %v", err)
	}

	newDir := filepath.Join(initialDir, "..")

	if _, err := executeCommand(core.CdCmd, newDir); err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	expected, err := filepath.Abs(newDir)
	if err != nil {
		t.Fatalf("Error getting the expected directory: %v", err)
	}

	got, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting the current directory: %v", err)
	}
	if got != expected {
		t.Errorf("Expected %q but got %q", expected, got)
	}

	if _, err := executeCommand(core.CdCmd, initialDir); err != nil {
		t.Fatalf("Failed to return to initial directory: %v", err)
	}
}

func TestMkdirCommand(t *testing.T) {
	dirName := "testdir"

	if _, err := executeCommand(core.MkdirCmd, dirName); err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		t.Fatalf("Expected directory %s to exist but it does not", dirName)
	}

	os.Remove(dirName)
}
