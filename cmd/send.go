package cmd

import (
	"github.com/fatih/color"
	"github.com/muhammedfazall/sendr-cli/internal/client"
	"github.com/muhammedfazall/sendr-cli/internal/config"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send an email",
	Run: func(cmd *cobra.Command, args []string) {
		to, _ := cmd.Flags().GetString("to")
		subject, _ := cmd.Flags().GetString("subject")
		body, _ := cmd.Flags().GetString("body")

		if to == "" || subject == "" || body == "" {
			fail("--to, --subject and --body are required")
		}

		cfg, err := config.Load()
		if err != nil {
			fail("Could not load config: %s", err)
		}

		if cfg.APIKey == "" {
			fail("No API key found. Run 'sendr keys create <name>' first")
		}

		c := client.New(cfg.APIURL, cfg.Token, cfg.APIKey)
		resp, err := c.SendEmail(to, subject, body)
		if err != nil {
			fail("Failed to send email: %s", err)
		}

		success("Email queued")
		color.HiBlack("  job_id: %s", resp.JobID)
	},
}

func init() {
	sendCmd.Flags().String("to", "", "Recipient email address")
	sendCmd.Flags().String("subject", "", "Email subject")
	sendCmd.Flags().String("body", "", "Email body")
}