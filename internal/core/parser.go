package core

import (
	"errors"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

type CommandChain struct {
	Commands  []Command
	Operators []string // "|", "&&", "||", ";"
}

func ParseInput(input string) (*CommandChain, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, nil
	}

	// Handle built-in exit command
	if input == "exit" {
		return &CommandChain{
			Commands: []Command{{Name: "exit", Args: []string{}}},
		}, nil
	}

	// Split by operators while preserving them
	tokens, operators := tokenize(input)
	
	if len(tokens) == 0 {
		return nil, errors.New("no commands found")
	}

	commands := make([]Command, 0, len(tokens))
	for _, token := range tokens {
		cmd, err := parseCommand(token)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}

	return &CommandChain{
		Commands:  commands,
		Operators: operators,
	}, nil
}

func tokenize(input string) ([]string, []string) {
	var tokens []string
	var operators []string
	var current strings.Builder
	
	i := 0
	for i < len(input) {
		switch {
		case i < len(input)-1 && input[i:i+2] == "&&":
			if current.Len() > 0 {
				tokens = append(tokens, strings.TrimSpace(current.String()))
				current.Reset()
			}
			operators = append(operators, "&&")
			i += 2
		case i < len(input)-1 && input[i:i+2] == "||":
			if current.Len() > 0 {
				tokens = append(tokens, strings.TrimSpace(current.String()))
				current.Reset()
			}
			operators = append(operators, "||")
			i += 2
		case input[i] == '|':
			if current.Len() > 0 {
				tokens = append(tokens, strings.TrimSpace(current.String()))
				current.Reset()
			}
			operators = append(operators, "|")
			i++
		case input[i] == ';':
			if current.Len() > 0 {
				tokens = append(tokens, strings.TrimSpace(current.String()))
				current.Reset()
			}
			operators = append(operators, ";")
			i++
		default:
			current.WriteByte(input[i])
			i++
		}
	}
	
	if current.Len() > 0 {
		tokens = append(tokens, strings.TrimSpace(current.String()))
	}
	
	return tokens, operators
}

func parseCommand(cmdStr string) (Command, error) {
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return Command{}, errors.New("empty command")
	}
	
	return Command{
		Name: parts[0],
		Args: parts[1:],
	}, nil
}
