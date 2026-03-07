package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wilblik/genie/internal/parser"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check [message]",
	Short: "Validate a commit message against Conventional Commits",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		msg := args[0]
		entry, err := parser.ParseCommitMessage(msg)
		if err != nil {
			return fmt.Errorf("error parsing message: %w", err)
		}

		if entry == nil {
			fmt.Println("❌ Message does not follow Conventional Commits standard.")
			fmt.Println("Example: feat(ui): add new button")
			return fmt.Errorf("invalid commit message")
		}

		fmt.Println("✅ Message is valid!")
		fmt.Printf("Type:    %s\n", entry.Type)

		if entry.Scope != "" {
			fmt.Printf("Scope:   %s\n", entry.Scope)
		}

		fmt.Printf("Subject: %s\n", entry.Subject)

		if entry.Body != "" {
			fmt.Println("\nBody:")
			fmt.Println(entry.Body)
		}

		if len(entry.Footers) > 0 {
			fmt.Println("\nFooters:")
			for k, v := range entry.Footers {
				fmt.Printf("%s: %s\n", k, v)
			}
		}

		if entry.IsBreaking {
			fmt.Println("\n⚠️  This is a BREAKING CHANGE")
		}

		return nil
	},
}
