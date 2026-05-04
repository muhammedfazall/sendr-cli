package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sendr",
	Short: "Sendr CLI — send transactional emails from your terminal",
	Long:  `Sendr CLI lets you send emails, manage API keys, and configure your Sendr account from the terminal.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Red("Error: %s", err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(keysCmd)
	rootCmd.AddCommand(configCmd)
}

// helpers used across commands

func success(format string, a ...any) {
	color.Green("✓ "+format, a...)
}

func fail(format string, a ...any) {
	color.Red("✗ "+format, a...)
	os.Exit(1)
}

func info(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}
