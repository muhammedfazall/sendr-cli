package cmd

import (
	"github.com/muhammedfazall/sendr-cli/internal/client"
	"github.com/muhammedfazall/sendr-cli/internal/config"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Sign out of your Sendr account",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fail("Could not load config: %s", err)
		}

		if cfg.Token != "" {
			c := client.New(cfg.APIURL, cfg.Token, cfg.APIKey)
			// best effort — don't fail if API call fails
			_ = c.Logout()
		}

		if err := cfg.Clear(); err != nil {
			fail("Could not clear config: %s", err)
		}

		success("Logged out")
	},
}