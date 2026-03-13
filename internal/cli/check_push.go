package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wilblik/genie/internal/config"
	"github.com/wilblik/genie/internal/git"
)

func init() {
	rootCmd.AddCommand(checkPushCmd)
}

var checkPushCmd = &cobra.Command{
	Use:    "check-push",
	Short:  "Internal: Validates a range of commits being pushed to a branch",
	Hidden: true,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read stdin from the pre-push hook
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Fields(line)
			if len(parts) < 4 { continue }

			remoteRef := parts[2]
			remoteSha := parts[3]
			localSha := parts[1]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("configuration file reading failed. Reason: %s\n\nPlease run 'genie init' first if %s is not present", err, config.ConfigFileName)
			}

			targetBranch := strings.TrimPrefix(remoteRef, "refs/heads/")
			if cfg.EnforceAll || targetBranch == cfg.ProtectedBranch {
				fmt.Printf("Validating commits for branch '%s'...\n", targetBranch)

				messages, err := git.GetCommitMessages(remoteSha, localSha)
				if err != nil { return err; }

				for _, msg := range messages {
					_, err := git.ParseCommitMessage(msg, cfg)
					if err != nil {
						fmt.Printf("❌ Invalid commit message in push range:\n%s\n", msg)
						return err
					}
				}
			}
		}

		return nil
	},
}
