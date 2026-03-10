package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wilblik/genie/internal/config"
)

const GithubWorkflowScriptTemplate = `name: Genie PR Title Check

on:
  pull_request:
    branches:
      - %s
    types: [opened, edited, synchronize, reopened]

jobs:
  check-pr-title:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Install Genie
        run: go install github.com/wilblik/genie/cmd/genie@latest

      - name: Validate PR Title
        run: genie check-msg "${{ github.event.pull_request.title }}"
`

func init() {
	initCmd.AddCommand(initGithubCmd)
}

var initGithubCmd = &cobra.Command{
	Use:   "github",
	Short: "Scaffold GitHub Actions workflow for PR title enforcement",
	Long:  `Creates a .github/workflows/genie.yml file to automatically validate Pull Request titles against your commit standard.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("configuration file not found. Please run 'genie init' first to set your protected branch")
		}

		workflowDir := filepath.Join(".github", "workflows")
		if err := os.MkdirAll(workflowDir, 0755); err != nil {
			return fmt.Errorf("failed to create workflows directory: %w", err)
		}

		workflowPath := filepath.Join(workflowDir, "genie.yml")

		if _, err := os.Stat(workflowPath); err == nil {
			fmt.Printf("⚠️  Workflow file %s already exists. Skipping creation.\n", workflowPath)
			return nil
		}

		script := fmt.Sprintf(GithubWorkflowScriptTemplate, cfg.ProtectedBranch)

		if err := os.WriteFile(workflowPath, []byte(script), 0644); err != nil {
			return fmt.Errorf("failed to write workflow file: %w", err)
		}

		fmt.Printf("✨ GitHub Action scaffolded successfully at %s\n\n", workflowPath)
		fmt.Println("GitHub Repository Setup Recipe:")
		fmt.Println("1. Go to your repository Settings > Branches.")
		fmt.Println("2. Add a branch protection rule for 'master' (or your main branch).")
		fmt.Println("3. Check 'Require status checks to pass before merging'.")
		fmt.Println("4. Search for 'check-pr-title' and make it required.")
		fmt.Println("5. Check 'Require linear history' or ensure Squash Merging is the only allowed merge strategy.")
		fmt.Println("\nThis ensures that all code landing in your protected branch has a perfectly formatted commit message!")

		return nil
	},
}
