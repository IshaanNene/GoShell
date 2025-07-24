package core

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Executor struct {
	lastExitCode int
}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) ExecuteChain(chain *CommandChain) error {
	if len(chain.Commands) == 0 {
		return nil
	}

	
	if len(chain.Commands) == 1 {
		return e.executeCommand(chain.Commands[0], nil, nil)
	}

	
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
			
			return e.executePipeChain(chain.Commands[i:])
		case "&&":
			
			err = e.executeCommand(cmd, nil, nil)
			if err != nil {
				return err
			}
		case "||":
			
			err = e.executeCommand(cmd, nil, nil)
			if err == nil {
				
				for j := i + 1; j < len(chain.Commands) && j < len(chain.Operators) && chain.Operators[j] == "||"; j++ {
					i = j
				}
			}
		case ";", "":
			
			err = e.executeCommand(cmd, nil, nil)
			
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
	
	
	for i, command := range commands {
		var execCmd *exec.Cmd
		var stdin io.Reader
		var stdout io.Writer = os.Stdout
		var stderr io.Writer = os.Stderr
		
		
		if e.isBuiltIn(command.Name) {
			return fmt.Errorf("built-in commands not fully supported in pipes: %s", command.Name)
		}
		
		execCmd = exec.Command(command.Name, command.Args...)
		
		
		if i == 0 {
			
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
			
			stdin = pipes[i-1]
		}
		execCmd.Stdin = stdin
		
		
		if i == len(commands)-1 {
			
			stdout, stderr = e.setupRedirection(command)
			if f, ok := stdout.(*os.File); ok && f != os.Stdout {
				defer f.Close()
			}
			if f, ok := stderr.(*os.File); ok && f != os.Stderr {
				defer f.Close()
			}
		} else {
			
			pr, pw := io.Pipe()
			pipes = append(pipes, pr)
			stdout = pw
		}
		execCmd.Stdout = stdout
		execCmd.Stderr = stderr
		
		cmds = append(cmds, execCmd)
	}
	
	
	for i, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start command %d: %v", i, err)
		}
	}
	
	
	go func() {
		for _, cmd := range cmds[:len(cmds)-1] {
			cmd.Wait()
			if pw, ok := cmd.Stdout.(io.WriteCloser); ok {
				pw.Close()
			}
		}
	}()
	
	
	return cmds[len(cmds)-1].Wait()
}

func (e *Executor) executeCommand(cmd Command, stdin io.Reader, stdout io.Writer) error {
	
	if e.isBuiltIn(cmd.Name) {
		return e.executeBuiltIn(cmd, stdin, stdout)
	}

	
	execCmd := exec.Command(cmd.Name, cmd.Args...)
	
	
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
	
	
	outWriter, errWriter := e.setupRedirection(cmd)
	if stdout != nil {
		execCmd.Stdout = stdout
	} else {
		execCmd.Stdout = outWriter
	}
	execCmd.Stderr = errWriter
	
	
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
	
	
	fmt.Fprintln(stdout, "History command - implementation would show command history")
	return nil
}