package cli

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"
	"github.com/wilblik/genie/internal/config"
	"github.com/wilblik/genie/internal/parser"
	"github.com/wilblik/genie/internal/models"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:          "check [message]",
	Short:        "Validate a commit message against Conventional Commits",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		msg := args[0]

		cfg, err := config.Load(config.ConfigFileName)
		if err != nil {
			return fmt.Errorf("configuration file not found. Please run 'genie init' first")
		}

		commitMsg, err := parser.ParseCommitMessage(msg)
		if err != nil {
			return fmt.Errorf("error parsing message: %w", err)
		}

		if err := validateMsg(cfg, commitMsg); err != nil {
			return err
		}

		printParsedMsg(commitMsg)

		return nil
	},
}

func validateMsg(cfg *config.Config, commitMsg *models.CommitMessage) error {
	if commitMsg == nil {
		return fmt.Errorf("message does not follow Conventional Commits standard\nExample: feat(ui): add new button")
	}

	if !slices.Contains(cfg.Types, commitMsg.Type) {
		return fmt.Errorf("type '%s' is not in the allowed types lsit: %v", commitMsg.Type, cfg.Types)
	}

	if cfg.RequireScope && commitMsg.Scope == "" {
		return fmt.Errorf("scope is required but was not provided")
	}

	if len(cfg.AllowedScopes) > 0 && commitMsg.Scope != "" {
		if !slices.Contains(cfg.AllowedScopes, commitMsg.Scope) {
			return fmt.Errorf("scope '%s' is not in the allowed scopes list: %v", commitMsg.Scope, cfg.AllowedScopes)
		}
	}

	return nil
}

func printParsedMsg(commitMsg *models.CommitMessage) {
	fmt.Println("✅ Message is valid!")
	fmt.Printf("Type:    %s\n", commitMsg.Type)
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
