package core

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Executor struct {
	// Add any executor state here if needed
}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) ExecuteChain(chain *CommandChain) error {
	if len(chain.Commands) == 0 {
		return nil
	}

	var lastErr error
	for i := 0; i < len(chain.Commands); i++ {
		cmd := chain.Commands[i]

		// Handle pipes within the chain
		if i < len(chain.Operators) && chain.Operators[i] == "|" {
			// Find the end of the pipe chain
			pipeEnd := i
			for pipeEnd < len(chain.Operators) && chain.Operators[pipeEnd] == "|" {
				pipeEnd++
			}

			pipelineCommands := chain.Commands[i : pipeEnd+1]
			lastErr = e.executePipeline(pipelineCommands)
			i = pipeEnd // Move index to the end of the pipeline
			continue
		}

		// Handle chaining operators (&&, ||, ;)
		if i > 0 && i-1 < len(chain.Operators) {
			op := chain.Operators[i-1]

			switch op {
			case "&&";
				if lastErr != nil {
					return lastErr // Stop on failure
				}
			case "||";
				if lastErr == nil {
					continue // Skip on success
				}
			case ";":
				// Continue regardless of error
			}
		}

		// Execute single command
		lastErr = e.executeCommand(cmd, nil, os.Stdout)
	}

	return lastErr
}

func (e *Executor) executePipeline(commands []Command) error {
	var execCmds []*exec.Cmd
	for _, cmd := range commands {
		execCmds = append(execCmds, e.createExecCommand(cmd))
	}

	// Connect pipes
	for i := 0; i < len(execCmds)-1; i++ {
		stdout, err := execCmds[i].StdoutPipe()
		if err != nil {
			return err
		}
		execCmds[i+1].Stdin = stdout
		execCmds[i].Stderr = os.Stderr // Redirect stderr to os.Stderr for all commands in pipe
	}
	execCmds[len(execCmds)-1].Stdout = os.Stdout
	execCmds[len(execCmds)-1].Stderr = os.Stderr

	// Start commands
	for _, cmd := range execCmds {
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	// Wait for commands to finish
	var lastErr error
	for _, cmd := range execCmds {
		if err := cmd.Wait(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func (e *Executor) executeCommand(cmd Command, stdin io.Reader, stdout io.Writer) error {
	// Handle built-in commands
	switch cmd.Name {
	case "exit":
		os.Exit(0)
	case "cd":
		return e.handleCD(cmd.Args)
	case "pwd":
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Fprintln(stdout, wd)
		return nil
	case "echo":
		fmt.Fprintln(stdout, strings.Join(cmd.Args, " "))
		return nil
	
	}

	// Execute external command
	execCmd := e.createExecCommand(cmd)
	if stdin != nil {
		execCmd.Stdin = stdin
	} else {
		execCmd.Stdin = os.Stdin
	}

	// Handle redirection for external commands
	// This is a basic implementation. A full shell would handle redirection in the parser and executor.
	
	
	execCmd.Stdout = stdout
	execCmd.Stderr = os.Stderr

	return execCmd.Run()
}

func (e *Executor) createExecCommand(cmd Command) *exec.Cmd {
	return exec.Command(cmd.Name, cmd.Args...)
}

func (e *Executor) handleCD(args []string) error {
	var dir string
	if len(args) == 0 {
		// cd with no arguments goes to home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		dir = homeDir
	} else {
		dir = args[0]
	}

	return os.Chdir(dir)
}
