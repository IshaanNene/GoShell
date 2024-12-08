package core

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"strings"
	"time"
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
