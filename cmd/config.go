package cmd

import (
	"github.com/fatih/color"
	"github.com/muhammedfazall/sendr-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

var configSetURLCmd = &cobra.Command{
	Use:   "set-url <url>",
	Short: "Set the Sendr API URL (for self-hosted instances)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fail("Could not load config: %s", err)
		}

		cfg.APIURL = args[0]
		if err := cfg.Save(); err != nil {
			fail("Could not save config: %s", err)
		}

		success("API URL set to %s", args[0])
	},
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fail("Could not load config: %s", err)
		}

		color.HiBlack("api_url: ") 
		color.HiWhite("  %s", cfg.APIURL)

		color.HiBlack("logged_in:")
		if cfg.Token != "" {
			color.Green("  yes")
		} else {
			color.Red("  no")
		}

		color.HiBlack("api_key:")
		if cfg.APIKey != "" {
			color.HiWhite("  %s...", cfg.APIKey[:20])
		} else {
			color.Red("  not set")
		}
	},
}

func init() {
	configCmd.AddCommand(configSetURLCmd)
	configCmd.AddCommand(configViewCmd)
}