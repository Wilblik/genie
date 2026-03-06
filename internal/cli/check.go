package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wilblik/genie/internal/parser"
)

var checkCmd = &cobra.Command{
	Use:   "check [message]",
	Short: "Validate a commit message against Conventional Commits",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		msg := args[0]
		entry, err := parser.ParseCommit(msg)
		if err != nil {
			fmt.Printf("❌ Error parsing message: %v\n", err)
			os.Exit(1)
		}

		if entry == nil {
			fmt.Println("❌ Message does not follow Conventional Commits standard.")
			fmt.Println("Example: feat(ui): add new button")
			os.Exit(1)
		}

		fmt.Println("✅ Message is valid!")
		fmt.Printf("Type:  %s\n", entry.Type)
		if entry.Scope != "" {
			fmt.Printf("Scope: %s\n", entry.Scope)
		}
		fmt.Printf("Subj:  %s\n", entry.Subject)
		if entry.IsBreaking {
			fmt.Println("⚠️  This is a BREAKING CHANGE")
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
