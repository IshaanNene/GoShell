package core

import (
	"fmt"
	"regexp"
	"strings"
)

type Command struct {
	Name        string
	Args        []string
	InputFile   string
	OutputFile  string
	AppendFile  string
	ErrorFile   string
}

type CommandChain struct {
	Commands  []Command
	Operators []string 
}


func ParseInput(input string) (*CommandChain, error) {
	if strings.TrimSpace(input) == "" {
		return nil, nil
	}

	
	parts, operators := splitByOperators(input)
	
	var commands []Command
	for _, part := range parts {
		cmd, err := parseCommand(strings.TrimSpace(part))
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


func splitByOperators(input string) ([]string, []string) {
	
	re := regexp.MustCompile(`(\|\||\&&|\||\;)`)
	
	
	parts := re.Split(input, -1)
	matches := re.FindAllString(input, -1)
	
	var cleanParts []string
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			cleanParts = append(cleanParts, trimmed)
		}
	}
	
	return cleanParts, matches
}


func parseCommand(cmdStr string) (Command, error) {
	cmd := Command{}
	
	
	cmdStr = handleRedirections(cmdStr, &cmd)
	
	
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return cmd, fmt.Errorf("empty command")
	}
	
	cmd.Name = parts[0]
	if len(parts) > 1 {
		cmd.Args = parts[1:]
	}
	
	return cmd, nil
}


func handleRedirections(cmdStr string, cmd *Command) string {
	
	if strings.Contains(cmdStr, "<") {
		parts := strings.Split(cmdStr, "<")
		if len(parts) == 2 {
			cmdStr = strings.TrimSpace(parts[0])
			inputFile := strings.TrimSpace(parts[1])
			
			if idx := strings.IndexAny(inputFile, ">"); idx != -1 {
				cmd.InputFile = strings.TrimSpace(inputFile[:idx])
				cmdStr += " " + inputFile[idx:]
			} else {
				cmd.InputFile = inputFile
			}
		}
	}
	
	
	if strings.Contains(cmdStr, ">>") {
		parts := strings.Split(cmdStr, ">>")
		if len(parts) == 2 {
			cmdStr = strings.TrimSpace(parts[0])
			cmd.AppendFile = strings.TrimSpace(parts[1])
		}
	} else if strings.Contains(cmdStr, ">") {
		
		parts := strings.Split(cmdStr, ">")
		if len(parts) == 2 {
			cmdStr = strings.TrimSpace(parts[0])
			outputFile := strings.TrimSpace(parts[1])
			
			
			if strings.HasSuffix(parts[0], "2") {
				cmdStr = strings.TrimSpace(parts[0][:len(parts[0])-1])
				cmd.ErrorFile = outputFile
			} else {
				cmd.OutputFile = outputFile
			}
		}
	}
	
	return cmdStr
}