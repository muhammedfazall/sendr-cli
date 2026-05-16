package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "sendr",
	Version: version,
	Short:   "Sendr CLI — send transactional emails from your terminal",
	Long:    `Sendr CLI lets you send emails, manage API keys, and configure your Sendr account from the terminal.`,
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

func truncateEnd(value string, maxLen int) string {
	if maxLen <= 0 || len(value) <= maxLen {
		return value
	}
	if maxLen <= 3 {
		return value[:maxLen]
	}
	return value[:maxLen-3] + "..."
}

func dateOnly(value string) string {
	if len(value) >= 10 {
		return value[:10]
	}
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}
