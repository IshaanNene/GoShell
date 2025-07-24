// internal/core/executor.go
package core

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Executor struct {
	lastExitCode int
	historyManager *HistoryManager // Add reference to history manager
}

func NewExecutor() *Executor {
	return &Executor{}
}

// SetHistoryManager allows setting the history manager for built-in history command
func (e *Executor) SetHistoryManager(hm *HistoryManager) {
	e.historyManager = hm
}

func (e *Executor) ExecuteChain(chain *CommandChain) error {
	if len(chain.Commands) == 0 {
		return nil
	}

	// Handle single command
	if len(chain.Commands) == 1 {
		return e.executeCommand(chain.Commands[0], nil, nil)
	}

	// Handle command chains
	return e.executeChainedCommands(chain)
}

func (e *Executor) executeChainedCommands(chain *CommandChain) error {
	for i, cmd := range chain.Commands {
		var operator string
		if i < len(chain.Operators) {
			operator = chain.Operators[i]
		}

		var err error
		switch operator {
		case "|":
			// Handle pipe - execute remaining commands as a pipe chain
			return e.executePipeChain(chain.Commands[i:])
		case "&&":
			// Execute only if previous succeeded
			err = e.executeCommand(cmd, nil, nil)
			if err != nil {
				return err
			}
		case "||":
			// Execute only if previous failed
			err = e.executeCommand(cmd, nil, nil)
			if err == nil {
				// Skip remaining OR commands
				for j := i + 1; j < len(chain.Commands) && j < len(chain.Operators) && chain.Operators[j] == "||"; j++ {
					i = j
				}
			}
		case ";", "":
			// Always execute
			err = e.executeCommand(cmd, nil, nil)
			// Continue regardless of error for semicolon
		}
		
		if operator != ";" && err != nil {
			e.lastExitCode = 1
		} else if err == nil {
			e.lastExitCode = 0
		}
	}

	return nil
}

func (e *Executor) executePipeChain(commands []Command) error {
	if len(commands) < 2 {
		return fmt.Errorf("pipe requires at least 2 commands")
	}

	var cmds []*exec.Cmd
	var pipes []io.ReadCloser
	var customCommands []bool // Track which commands are custom
	
	// Create all commands
	for i, command := range commands {
		var execCmd *exec.Cmd
		var stdin io.Reader
		var stdout io.Writer = os.Stdout
		var stderr io.Writer = os.Stderr
		isCustom := false
		
		// Check if it's a registered command first
		if IsCommandRegistered(command.Name) {
			// For pipes with custom commands, we need special handling
			// Custom commands in pipes are complex, so we'll limit this for now
			return fmt.Errorf("custom registered commands not fully supported in pipes yet: %s", command.Name)
		}
		
		// Handle built-ins for pipes (limited support)
		if e.isBuiltIn(command.Name) {
			return fmt.Errorf("built-in commands not fully supported in pipes: %s", command.Name)
		}
		
		execCmd = exec.Command(command.Name, command.Args...)
		customCommands = append(customCommands, isCustom)
		
		// Set up stdin
		if i == 0 {
			// First command - handle input redirection
			if command.InputFile != "" {
				file, err := os.Open(command.InputFile)
				if err != nil {
					return fmt.Errorf("cannot open input file %s: %v", command.InputFile, err)
				}
				defer file.Close()
				stdin = file
			} else {
				stdin = os.Stdin
			}
		} else {
			// Use output from previous command
			stdin = pipes[i-1]
		}
		execCmd.Stdin = stdin
		
		// Set up stdout
		if i == len(commands)-1 {
			// Last command - handle output redirection
			stdout, stderr = e.setupRedirection(command)
			if f, ok := stdout.(*os.File); ok && f != os.Stdout {
				defer f.Close()
			}
			if f, ok := stderr.(*os.File); ok && f != os.Stderr {
				defer f.Close()
			}
		} else {
			// Create pipe to next command
			pr, pw := io.Pipe()
			pipes = append(pipes, pr)
			stdout = pw
		}
		execCmd.Stdout = stdout
		execCmd.Stderr = stderr
		
		cmds = append(cmds, execCmd)
	}
	
	// Start all commands
	for i, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start command %d: %v", i, err)
		}
	}
	
	// Close write ends of pipes and wait for commands
	go func() {
		for _, cmd := range cmds[:len(cmds)-1] {
			cmd.Wait()
			if pw, ok := cmd.Stdout.(io.WriteCloser); ok {
				pw.Close()
			}
		}
	}()
	
	// Wait for the last command
	return cmds[len(cmds)-1].Wait()
}

func (e *Executor) executeCommand(cmd Command, stdin io.Reader, stdout io.Writer) error {
	// 1. First priority: Check for registered custom commands
	if IsCommandRegistered(cmd.Name) {
		return e.executeRegisteredCommand(cmd, stdin, stdout)
	}
	
	// 2. Second priority: Handle built-in commands
	if e.isBuiltIn(cmd.Name) {
		return e.executeBuiltIn(cmd, stdin, stdout)
	}

	// 3. Third priority: Execute external command
	return e.executeExternalCommand(cmd, stdin, stdout)
}

func (e *Executor) executeRegisteredCommand(cmd Command, stdin io.Reader, stdout io.Writer) error {
	cobraCmd := GetCommand(cmd.Name)
	if cobraCmd == nil {
		return fmt.Errorf("registered command not found: %s", cmd.Name)
	}
	
	// Create a copy of the command to avoid modifying the original
	cmdCopy := &cobra.Command{
		Use:                   cobraCmd.Use,
		Short:                 cobraCmd.Short,
		Long:                  cobraCmd.Long,
		Run:                   cobraCmd.Run,
		Args:                  cobraCmd.Args,
	}
	
	// Copy flags from original command
	cobraCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		cmdCopy.Flags().AddFlag(flag)
	})
	
	// Set up I/O redirection for custom commands
	outWriter, errWriter := e.setupRedirection(cmd)
	
	// Set output destinations
	if stdout != nil {
		cmdCopy.SetOut(stdout)
	} else {
		cmdCopy.SetOut(outWriter)
	}
	cmdCopy.SetErr(errWriter)
	
	// Handle input redirection
	if cmd.InputFile != "" {
		file, err := os.Open(cmd.InputFile)
		if err != nil {
			return fmt.Errorf("cannot open input file %s: %v", cmd.InputFile, err)
		}
		defer file.Close()
		cmdCopy.SetIn(file)
	} else if stdin != nil {
		cmdCopy.SetIn(stdin)
	} else {
		cmdCopy.SetIn(os.Stdin)
	}
	
	// Temporarily redirect os.Stdout and os.Stderr for commands that don't use cobra's writers
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	
	if stdout != nil {
		if f, ok := stdout.(*os.File); ok {
			os.Stdout = f
		}
	} else if outWriter != os.Stdout {
		if f, ok := outWriter.(*os.File); ok {
			os.Stdout = f
		}
	}
	
	if errWriter != os.Stderr {
		if f, ok := errWriter.(*os.File); ok {
			os.Stderr = f
		}
	}
	
	// Execute the cobra command
	cmdCopy.SetArgs(cmd.Args)
	err := cmdCopy.Execute()
	
	// Restore original stdout and stderr
	os.Stdout = originalStdout
	os.Stderr = originalStderr
	
	// Close files if they were opened for redirection
	if f, ok := outWriter.(*os.File); ok && f != os.Stdout && f != originalStdout {
		f.Close()
	}
	if f, ok := errWriter.(*os.File); ok && f != os.Stderr && f != originalStderr {
		f.Close()
	}
	
	return err
}

func (e *Executor) executeExternalCommand(cmd Command, stdin io.Reader, stdout io.Writer) error {
	// Create external command
	execCmd := exec.Command(cmd.Name, cmd.Args...)
	
	// Set up I/O
	if stdin != nil {
		execCmd.Stdin = stdin
	} else if cmd.InputFile != "" {
		file, err := os.Open(cmd.InputFile)
		if err != nil {
			return fmt.Errorf("cannot open input file %s: %v", cmd.InputFile, err)
		}
		defer file.Close()
		execCmd.Stdin = file
	} else {
		execCmd.Stdin = os.Stdin
	}
	
	// Set up output and error redirection
	outWriter, errWriter := e.setupRedirection(cmd)
	if stdout != nil {
		execCmd.Stdout = stdout
	} else {
		execCmd.Stdout = outWriter
	}
	execCmd.Stderr = errWriter
	
	// Close files after execution
	if f, ok := outWriter.(*os.File); ok && f != os.Stdout {
		defer f.Close()
	}
	if f, ok := errWriter.(*os.File); ok && f != os.Stderr {
		defer f.Close()
	}
	
	return execCmd.Run()
}

func (e *Executor) setupRedirection(cmd Command) (io.Writer, io.Writer) {
	var stdout io.Writer = os.Stdout
	var stderr io.Writer = os.Stderr
	
	// Handle output redirection
	if cmd.OutputFile != "" {
		file, err := os.Create(cmd.OutputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot create output file %s: %v\n", cmd.OutputFile, err)
		} else {
			stdout = file
		}
	} else if cmd.AppendFile != "" {
		file, err := os.OpenFile(cmd.AppendFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot open append file %s: %v\n", cmd.AppendFile, err)
		} else {
			stdout = file
		}
	}
	
	// Handle error redirection
	if cmd.ErrorFile != "" {
		file, err := os.Create(cmd.ErrorFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot create error file %s: %v\n", cmd.ErrorFile, err)
		} else {
			stderr = file
		}
	}
	
	return stdout, stderr
}

func (e *Executor) isBuiltIn(name string) bool {
	builtIns := []string{"exit", "cd", "pwd", "echo", "history"}
	for _, builtin := range builtIns {
		if name == builtin {
			return true
		}
	}
	return false
}

func (e *Executor) executeBuiltIn(cmd Command, stdin io.Reader, stdout io.Writer) error {
	if stdout == nil {
		outWriter, _ := e.setupRedirection(cmd)
		stdout = outWriter
		if f, ok := stdout.(*os.File); ok && f != os.Stdout {
			defer f.Close()
		}
	}
	
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
	case "history":
		return e.handleHistory(cmd.Args, stdout)
	}
	
	return fmt.Errorf("unknown built-in command: %s", cmd.Name)
}

func (e *Executor) handleCD(args []string) error {
	var dir string
	if len(args) == 0 {
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

func (e *Executor) handleHistory(args []string, stdout io.Writer) error {
	if e.historyManager == nil {
		fmt.Fprintln(stdout, "History not available")
		return nil
	}
	
	if len(args) == 0 {
		// Show all history
		history := e.historyManager.GetHistory()
		for i, cmd := range history {
			fmt.Fprintf(stdout, "%4d  %s\n", i+1, cmd)
		}
	} else {
		// Handle history arguments (like history 10)
		e.historyManager.PrintHistory(args)
	}
	
	return nil
}

// GetLastExitCode returns the exit code of the last executed command
func (e *Executor) GetLastExitCode() int {
	return e.lastExitCode
}