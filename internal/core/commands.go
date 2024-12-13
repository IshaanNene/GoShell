package core

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"strings"
	"time"
	"io"
	"bufio"
)

var PwdCmd = &cobra.Command{
	Use:   "pwd",
	Short: "Print the current working directory",
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		checkError(err, "getting current working directory")
		fmt.Println(dir)
	},
}

var DateCmd = &cobra.Command{
	Use:   "date",
	Short: "Print the current date and time",
	Run: func(cmd *cobra.Command, args []string) {
		today := time.Now()
		zone, _ := today.Zone()
		fmt.Printf("%s %s %d %02d:%02d:%02d %s %d\n",
			today.Weekday().String()[0:3],
			today.Month().String()[0:3],
			today.Day(),
			today.Hour(),
			today.Minute(),
			today.Second(),
			zone,
			today.Year())
	},
}

var Iamwho = &cobra.Command{
	Use:   "iamwho",
	Short: "The whoami command",
	Run: func(cmd *cobra.Command, args []string) {
		g, err := user.Current()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(g.Username)
		}
	},
}

var EchoCmd = &cobra.Command{
	Use:   "echo [text]",
	Short: "Display a line of text",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(strings.Join(args, " "))
	},
}

var ClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the terminal screen",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("\033[H\033[2J")
	},
}

var HeadCmd = &cobra.Command{
	Use:   "head [file]",
	Short: "Display the first 10 lines of a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(args[0])
		checkError(err, "opening file")
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for i := 0; i < 10 && scanner.Scan(); i++ {
			fmt.Println(scanner.Text())
		}
		checkError(scanner.Err(), "reading file")
	},
}

var TailCmd = &cobra.Command{
	Use:   "tail [file]",
	Short: "Display the last 10 lines of a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(args[0])
		checkError(err, "opening file")
		defer file.Close()

		lines, err := io.ReadAll(file)
		checkError(err, "reading file")
		allLines := strings.Split(string(lines), "\n")
		start := len(allLines) - 10
		if start < 0 {
			start = 0
		}
		for _, line := range allLines[start:] {
			fmt.Println(line)
		}
	},
}
