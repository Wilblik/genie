package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wilblik/genie/internal/config"
	"github.com/wilblik/genie/internal/git"
	"github.com/wilblik/genie/internal/models"
)

func init() {
	rootCmd.AddCommand(checkMsgCmd)
}

var checkMsgCmd = &cobra.Command{
	Use:          "check-msg [message]",
	Short:        "Validate a commit message against Conventional Commits",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		msg := args[0]

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("configuration file reading failed. Please run 'genie init' first if %s is not present", config.ConfigFileName)
		}

		commitMsg, err := git.ParseCommitMessage(msg, cfg)
		if err != nil { return err }

		printParsedMsg(commitMsg)

		return nil
	},
}

func printParsedMsg(commitMsg *models.CommitMessage) {
	fmt.Println("✅ Message is valid!")
	fmt.Printf("Type:    %s\n", commitMsg.ChangeType)
	if commitMsg.Scope != "" {
		fmt.Printf("Scope:   %s\n", commitMsg.Scope)
	}
	fmt.Printf("Subject: %s\n", commitMsg.Subject)

	if commitMsg.Body != "" {
		fmt.Println("\nBody:")
		fmt.Println(commitMsg.Body)
	}

	if len(commitMsg.Footers) > 0 {
		fmt.Println("\nFooters:")
		for k, v := range commitMsg.Footers {
			fmt.Printf("- %s: %s\n", k, v)
		}
	}

	if commitMsg.IsBreaking {
		fmt.Println("\n⚠️  This is a BREAKING CHANGE")
	}
}
