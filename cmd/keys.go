package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/muhammedfazall/sendr-cli/internal/client"
	"github.com/muhammedfazall/sendr-cli/internal/config"
	"github.com/spf13/cobra"
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage your API keys",
}

var keysCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new API key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fail("Could not load config: %s", err)
		}

		if cfg.Token == "" {
			fail("Not logged in. Run 'sendr login' first")
		}

		c := client.New(cfg.APIURL, cfg.Token, cfg.APIKey)
		key, err := c.CreateKey(args[0])
		if err != nil {
			fail("Could not create key: %s", err)
		}

		// Save the full key locally for use in send
		cfg.APIKey = key.APIKey
		if err := cfg.Save(); err != nil {
			fail("Could not save config: %s", err)
		}

		success("API key created and saved")
		fmt.Println()
		color.Yellow("  ⚠  This is shown only once — it is now saved to your config")
		color.HiWhite("  %s", key.APIKey)
		fmt.Println()
		color.HiBlack("  name:   %s", key.Name)
		color.HiBlack("  prefix: mk_live_%s...", key.Prefix)
	},
}

var keysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your API keys",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fail("Could not load config: %s", err)
		}

		if cfg.Token == "" {
			fail("Not logged in. Run 'sendr login' first")
		}

		c := client.New(cfg.APIURL, cfg.Token, cfg.APIKey)
		keys, err := c.ListKeys()
		if err != nil {
			fail("Could not list keys: %s", err)
		}

		if len(keys) == 0 {
			info("No API keys found")
			return
		}

		fmt.Printf("\n%-36s  %-20s  %-18s  %s\n",
			color.HiBlackString("ID"),
			color.HiBlackString("NAME"),
			color.HiBlackString("PREFIX"),
			color.HiBlackString("CREATED"),
		)
		fmt.Println(color.HiBlackString("─────────────────────────────────────────────────────────────────────────────────────"))

		for _, k := range keys {
			fmt.Printf("%-36s  %-20s  %-18s  %s\n",
				k.ID,
				k.Name,
				"mk_live_"+k.Prefix+"...",
				dateOnly(k.CreatedAt),
			)
		}
		fmt.Println()
	},
}

var keysRevokeCmd = &cobra.Command{
	Use:   "revoke <id>",
	Short: "Revoke an API key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fail("Could not load config: %s", err)
		}

		if cfg.Token == "" {
			fail("Not logged in. Run 'sendr login' first")
		}

		c := client.New(cfg.APIURL, cfg.Token, cfg.APIKey)
		if err := c.RevokeKey(args[0]); err != nil {
			fail("Could not revoke key: %s", err)
		}

		success("Key revoked")
	},
}

func init() {
	keysCmd.AddCommand(keysCreateCmd)
	keysCmd.AddCommand(keysListCmd)
	keysCmd.AddCommand(keysRevokeCmd)
}
