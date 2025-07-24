package core

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var CatCmd = &cobra.Command{
	Use:   "cat",
	Short: "Concatenate and display file contents",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		number, _ := cmd.Flags().GetBool("number")
		numberNonBlank, _ := cmd.Flags().GetBool("number-nonblank")
		squeezeBlank, _ := cmd.Flags().GetBool("squeeze-blank")
		showEnds, _ := cmd.Flags().GetBool("show-ends")
		showTabs, _ := cmd.Flags().GetBool("show-tabs")
		showNonPrinting, _ := cmd.Flags().GetBool("show-nonprinting")

		
		if len(args) >= 3 && args[len(args)-2] == "mx" {
			mergeFiles(args[:len(args)-2], args[len(args)-1])
		} else {
			
			if len(args) == 0 {
				displayStdin(number, numberNonBlank, squeezeBlank, showEnds, showTabs, showNonPrinting)
			} else {
				for _, file := range args {
					displayFile(file, number, numberNonBlank, squeezeBlank, showEnds, showTabs, showNonPrinting)
				}
			}
		}
	},
}

func init() {
	CatCmd.Flags().BoolP("number", "n", false, "Number all output lines")
	CatCmd.Flags().BoolP("number-nonblank", "b", false, "Number non-empty output lines, overrides -n")
	CatCmd.Flags().BoolP("squeeze-blank", "s", false, "Suppress repeated empty output lines")
	CatCmd.Flags().BoolP("show-ends", "E", false, "Display $ at end of each line")
	CatCmd.Flags().BoolP("show-tabs", "T", false, "Display TAB characters as ^I")
	CatCmd.Flags().BoolP("show-nonprinting", "v", false, "Use ^ and M- notation, except for LFD and TAB")
}

func mergeFiles(files []string, outputFile string) {
	var container []byte
	for _, file := range files {
		contents, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Error reading file %s: %v", file, err)
		}
		container = append(container, contents...)
	}

	err := os.WriteFile(outputFile, container, 0644)
	if err != nil {
		log.Fatalf("Error writing to file %s: %v", outputFile, err)
	}
	fmt.Println("Contents successfully written to", outputFile)
}

func displayStdin(number, numberNonBlank, squeezeBlank, showEnds, showTabs, showNonPrinting bool) {
	scanner := bufio.NewScanner(os.Stdin)

	lineNumber := 1
	prevLineEmpty := false
	for scanner.Scan() {
		line := scanner.Text()
		processAndPrintLine(line, &lineNumber, &prevLineEmpty, number, numberNonBlank, squeezeBlank, showEnds, showTabs, showNonPrinting)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from stdin: %v", err)
	}
}

func displayFile(file string, number, numberNonBlank, squeezeBlank, showEnds, showTabs, showNonPrinting bool) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", file, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	lineNumber := 1
	prevLineEmpty := false
	for scanner.Scan() {
		line := scanner.Text()
		processAndPrintLine(line, &lineNumber, &prevLineEmpty, number, numberNonBlank, squeezeBlank, showEnds, showTabs, showNonPrinting)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning file %s: %v", file, err)
	}
}

func processAndPrintLine(line string, lineNumber *int, prevLineEmpty *bool, number, numberNonBlank, squeezeBlank, showEnds, showTabs, showNonPrinting bool) {
	
	if squeezeBlank && (*prevLineEmpty) && line == "" {
		return
	}

	
	processedLine := line

	if showTabs == true {
		processedLine = strings.ReplaceAll(processedLine, "\t", "^I")
	}

	if showNonPrinting == true {
		processedLine = showNonPrintableChars(processedLine)
	}

	if showEnds == true {
		processedLine += "$"
	}

	
	if numberNonBlank == true && line != "" {
		fmt.Printf("%6d\t%s\n", *lineNumber, processedLine)
		(*lineNumber)++
	} else if number == true {
		fmt.Printf("%6d\t%s\n", *lineNumber, processedLine)
		(*lineNumber)++
	} else {
		fmt.Println(processedLine)
	}

	*prevLineEmpty = (line == "")
}

func showNonPrintableChars(line string) string {
	var result strings.Builder
	for _, r := range line {
		if r < 32 && r != '\t' && r != '\n' {
			
			if r == 127 {
				result.WriteString("^?")
			} else {
				result.WriteString(fmt.Sprintf("^%c", r+64))
			}
		} else if r == 127 {
			result.WriteString("^?")
		} else if r > 127 {
			
			if r < 160 {
				result.WriteString(fmt.Sprintf("M-^%c", r-64))
			} else {
				result.WriteString(fmt.Sprintf("M-%c", r-128))
			}
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
