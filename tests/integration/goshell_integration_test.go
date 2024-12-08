package integration

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCommandExecution(t *testing.T) {
	binaryPath := "./goshell"

	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		postCheck      func(t *testing.T)
		cleanup        func()
	}{
		{
			name:           "Test 'ls' command",
			args:           []string{"ls"},
			expectedOutput: "goshell\n",
		},
		{
			name: "Test 'pwd' command",
			args: []string{"pwd"},
			expectedOutput: func() string {
				dir, _ := os.Getwd()
				return dir + "\n"
			}(),
		},
		{
			name: "Test 'touch' command",
			args: []string{"touch", "testfile.txt"},
			postCheck: func(t *testing.T) {
				if _, err := os.Stat("testfile.txt"); os.IsNotExist(err) {
					t.Fatalf("Expected file 'testfile.txt' to exist but it does not")
				}
			},
			cleanup: func() {
				os.Remove("testfile.txt")
			},
		},
		{
			name: "Test 'cd' command to home directory",
			args: []string{"cd", "~"},
			postCheck: func(t *testing.T) {
				usr, _ := os.UserHomeDir()
				dir, _ := os.Getwd()
				if dir != usr {
					t.Fatalf("Expected to be in home directory but got %s", dir)
				}
			},
		},
		{
			name: "Test 'cd' command to parent directory",
			args: []string{"cd", ".."},
			postCheck: func(t *testing.T) {
				currentDir, _ := os.Getwd()
				parentDir := filepath.Dir(currentDir)
				dir, _ := os.Getwd()
				if dir != parentDir {
					t.Fatalf("Expected to be in parent directory but got %s", dir)
				}
			},
		},
		{
			name: "Test 'mkdir' command",
			args: []string{"mkdir", "testdir"},
			postCheck: func(t *testing.T) {
				if _, err := os.Stat("testdir"); os.IsNotExist(err) {
					t.Fatalf("Expected directory 'testdir' to exist but it does not")
				}
			},
			cleanup: func() {
				os.Remove("testdir")
			},
		},
		{
			name: "Test 'rm' command",
			args: []string{"rm", "testfile.txt"},
			postCheck: func(t *testing.T) {
				if _, err := os.Stat("testfile.txt"); !os.IsNotExist(err) {
					t.Fatalf("Expected file 'testfile.txt' to be removed but it still exists")
				}
			},
			cleanup: func() {
				os.Remove("testfile.txt")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out

			err := cmd.Run()
			if err != nil {
				t.Fatalf("Error executing command: %v", err)
			}

			if tt.expectedOutput != "" {
				if got := out.String(); got != tt.expectedOutput {
					t.Errorf("Expected %q but got %q", tt.expectedOutput, got)
				}
			}

			if tt.postCheck != nil {
				tt.postCheck(t)
			}

			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}
