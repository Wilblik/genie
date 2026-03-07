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

			cfg, err := config.Load(config.ConfigFileName)
			if err != nil { return err; }

			targetBranch := strings.TrimPrefix(remoteRef, "refs/heads/")
			if cfg.EnforceAll || targetBranch == cfg.ProtectedBranch {
				fmt.Printf("Genie: Validating commits for branch '%s'...\n", targetBranch)

				messages, err := git.GetCommitMessages(remoteSha, localSha)
				if err != nil { return err; }

				for _, msg := range messages {
					commitMsg, err := git.ParseCommitMessage(msg)
					if err != nil { return err; }

					if err := git.ValidateCommitMessage(cfg, commitMsg); err != nil {
						fmt.Printf("❌ Invalid commit message in push range:\n%s\n", msg)
						return err
					}
				}
			}
		}

		return nil
	},
}
