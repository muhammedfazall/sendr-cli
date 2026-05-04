package cmd

import (
	"github.com/muhammedfazall/sendr-cli/internal/auth"
	"github.com/muhammedfazall/sendr-cli/internal/config"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with your Sendr account",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fail("Could not load config: %s", err)
		}

		token, err := auth.Login(cfg.APIURL)
		if err != nil {
			fail("Login failed: %s", err)
		}

		cfg.Token = token
		if err := cfg.Save(); err != nil {
			fail("Could not save config: %s", err)
		}

		success("Logged in successfully")
	},
}