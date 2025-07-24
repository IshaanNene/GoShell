package core

import (
	"fmt"
	"strings"
	"text/scanner"
)

// ParsedCommand represents a single command with its arguments and potential redirection.
type ParsedCommand struct {
	Name string
	Args []string
	RedirectInput string
	RedirectOutput string
	AppendOutput bool
}

// CommandPipeline represents a sequence of commands connected by pipes.
type CommandPipeline struct {
	Commands []*ParsedCommand
}

// CommandChain represents a sequence of command pipelines connected by chaining operators (&&, ||, ;).
type CommandChain struct {
	Pipelines []*CommandPipeline
	Operators []string // &&, ||, or ;
}

// ParseInput parses the input string into a CommandChain.
func ParseInput(input string) (*CommandChain, error) {
	// Trim leading/trailing whitespace
	input = strings.TrimSpace(input)

	// Handle empty input
	if input == "" {
		return nil, nil
	}

	chain := &CommandChain{}

	// Split by chaining operators (&&, ||, ;)
	chainOperators := []string{"&&", "||", ";"}
	pipelineStrs, operators := splitByOperators(input, chainOperators)

	for _, pipelineStr := range pipelineStrs {
		pipelineStr = strings.TrimSpace(pipelineStr)
		if pipelineStr == "" {
			continue
		}

		pipeline := &CommandPipeline{}
		// Split by pipe operator
		commandStrs := strings.Split(pipelineStr, "|")

		for _, commandStr := range commandStrs {
			commandStr = strings.TrimSpace(commandStr)
			if commandStr == "" {
				return nil, fmt.Errorf("syntax error: empty command in pipeline")
			}

			parsedCmd, err := parseSingleCommand(commandStr)
			if err != nil {
				return nil, err
			}

			pipeline.Commands = append(pipeline.Commands, parsedCmd)
		}

		chain.Pipelines = append(chain.Pipelines, pipeline)
	}

	chain.Operators = operators

	return chain, nil
}

// parseSingleCommand parses a single command string, including redirection.
func parseSingleCommand(commandStr string) (*ParsedCommand, error) {
	parsedCmd := &ParsedCommand{}

	var s scanner.Scanner
	s.Init(strings.NewReader(commandStr))
	s.Mode = scanner.ScanWords | scanner.ScanStrings

	// Read the command name
	token := s.Scan()
	if token == scanner.EOF {
		return nil, fmt.Errorf("syntax error: empty command")
	}
	parsedCmd.Name = s.TokenText()

	// Read arguments and redirection operators
	args := []string{}
	for token != scanner.EOF {
		token = s.Scan()
		part := s.TokenText()

		switch part {
		case ">";
			token = s.Scan()
			if token == scanner.EOF || s.TokenText() == "" {
				return nil, fmt.Errorf("syntax error: no output file specified after >")
			}
			parsedCmd.RedirectOutput = s.TokenText()
			parsedCmd.AppendOutput = false
		case ">>":
			token = s.Scan()
			if token == scanner.EOF || s.TokenText() == "" {
				return nil, fmt.Errorf("syntax error: no output file specified after >>")
			}
			parsedCmd.RedirectOutput = s.TokenText()
			parsedCmd.AppendOutput = true
		case "<":
			token = s.Scan()
			if token == scanner.EOF || s.TokenText() == "" {
				return nil, fmt.Errorf("syntax error: no input file specified after <")
			}
			parsedCmd.RedirectInput = s.TokenText()
		default:
			if token != scanner.EOF {
				args = append(args, part)
			}
		}
	}

	parsedCmd.Args = args

	return parsedCmd, nil
}

// splitByOperators splits a string by a list of operators, returning the parts and the operators found.
func splitByOperators(input string, operators []string) ([]string, []string) {
	var parts []string
	var ops []string
	lastIndex := 0

	for i := 0; i < len(input); {
		foundOperator := false
		for _, op := range operators {
			if strings.HasPrefix(input[i:], op) {
				parts = append(parts, input[lastIndex:i])
				ops = append(ops, op)
				lastIndex = i + len(op)
				i += len(op)
				foundOperator = true
				break
			}
		}
		if !foundOperator {
			i++
		}
	}

	parts = append(parts, input[lastIndex:])

	return parts, ops
}
