package main

import (
	"fmt"
	"jonathantyar/tokopedia-crawler/cmd"
	_ "jonathantyar/tokopedia-crawler/src/config"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Run Command",
}

func init() {
	rootCmd.AddCommand(cmd.ScrapperCmd)
	rootCmd.AddCommand(cmd.GooseCmd)
}

func main() {
	defer func() {
		e := recover()
		if e != nil {
			fmt.Println(e)
			fmt.Println(string(debug.Stack()))
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
