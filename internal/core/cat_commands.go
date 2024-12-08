package core

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	_ "github.com/thanhpk/ascii"
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

		if len(args) > 1 && args[len(args)-2] == "mx" {
			mergeFiles(args[:len(args)-2], args[len(args)-1])
		} else {
			for _, file := range args {
				displayFile(file, number, numberNonBlank, squeezeBlank, showEnds, showTabs)
			}
		}
	},
}

func init() {
	// Define the flags
	CatCmd.Flags().BoolP("number", "n", false, "Number all output lines")
	CatCmd.Flags().BoolP("number-nonblank", "b", false, "Number non-empty output lines, overrides -n")
	CatCmd.Flags().BoolP("squeeze-blank", "s", false, "Suppress repeated empty output lines")
	CatCmd.Flags().BoolP("show-ends", "E", false, "Display $ at end of each line")
	CatCmd.Flags().BoolP("show-tabs", "T", false, "Display TAB characters as ^I")
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

func displayFile(file string, number, numberNonBlank, squeezeBlank, showEnds, showTabs bool) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", file, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	fmt.Println("Contents of:", file)

	lineNumber := 1
	prevLineEmpty := false
	for scanner.Scan() {
		line := scanner.Text()
		if squeezeBlank && prevLineEmpty && line == "" {
			continue
		}
		if showTabs {
			line = strings.ReplaceAll(line, "\t", "^I")
		}
		if showEnds {
			line += "$"
		}
		if numberNonBlank && line != "" {
			fmt.Printf("    %d %s\n", lineNumber, line)
			lineNumber++
		} else if number {
			fmt.Printf("    %d %s\n", lineNumber, line)
			lineNumber++
		} else {
			fmt.Println(line)
		}
		prevLineEmpty = (line == "")
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning file %s: %v", file, err)
	}
	fmt.Println()
}
