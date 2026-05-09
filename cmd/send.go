package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/muhammedfazall/sendr-cli/internal/client"
	"github.com/muhammedfazall/sendr-cli/internal/config"
	"github.com/spf13/cobra"
)

// spinner frames for the polling animation
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send an email",
	Run: func(cmd *cobra.Command, args []string) {
		to, _ := cmd.Flags().GetString("to")
		subject, _ := cmd.Flags().GetString("subject")
		bodyFlag, _ := cmd.Flags().GetString("body")
		noWait, _ := cmd.Flags().GetBool("no-wait")

		if to == "" || subject == "" || bodyFlag == "" {
			fail("--to, --subject and --body are required")
		}

		// Step 3: resolve body — if starts with @, read from file
		body, err := resolveBody(bodyFlag)
		if err != nil {
			fail("Could not read body file: %s", err)
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

		// Step 2: poll until terminal state unless --no-wait
		if noWait {
			return
		}

		fmt.Println()
		pollUntilDone(c, resp.JobID)
	},
}

// resolveBody returns the body string.
// If value starts with @, it reads from the file path that follows.
func resolveBody(value string) (string, error) {
	if !strings.HasPrefix(value, "@") {
		return value, nil
	}
	path := strings.TrimPrefix(value, "@")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// pollUntilDone polls GET /emails/:id every 2s, showing a spinner,
// until the job reaches a terminal state (sent or failed).
func pollUntilDone(c *client.Client, jobID string) {
	frame := 0
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Print initial spinner line
	fmt.Printf("  %s Waiting for delivery...", spinnerFrames[frame])

	for range ticker.C {
		frame = (frame + 1) % len(spinnerFrames)

		job, err := c.GetJob(jobID)
		if err != nil {
			// Clear line and warn — keep polling, transient errors happen
			fmt.Printf("\r  %s Waiting for delivery...         ", spinnerFrames[frame])
			continue
		}

		if !job.Done() {
			fmt.Printf("\r  %s Waiting for delivery... (%s)", spinnerFrames[frame], job.Status)
			continue
		}

		// Clear spinner line before printing final result
		fmt.Print("\r" + strings.Repeat(" ", 60) + "\r")

		switch job.Status {
		case "sent":
			success("Delivered")
		case "failed":
			color.Red("✗ Delivery failed (retries: %d)", job.Retries)
		}
		return
	}
}

func init() {
	sendCmd.Flags().String("to", "", "Recipient email address")
	sendCmd.Flags().String("subject", "", "Email subject")
	sendCmd.Flags().String("body", "", "Email body, or @path/to/file.txt to read from file")
	sendCmd.Flags().Bool("no-wait", false, "Return immediately after queuing without waiting for delivery")
}