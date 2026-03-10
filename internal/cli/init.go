package cli

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/wilblik/genie/internal/config"
	"github.com/wilblik/genie/internal/git"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive setup for Genie in your repository",
	Long:  `Guides you through setting up Genie, including protected branches, scopes, and enforcement levels.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(config.ConfigFileName); err == nil {
			if !promptOverwrite() {
				fmt.Println("Init cancelled.")
				return nil;
			}
		}

		cfg := config.NewDefaultConfig()
		if err := promptBranch(cfg);  err != nil { return err }
		if err := promptScope(cfg);   err != nil { return err }
		if err := promptEnforce(cfg); err != nil { return err }
		if err := promptTypes(cfg);   err != nil { return err }
		if err := promptScopes(cfg);  err != nil { return err }
		if err := promptGithubAction(cfg); err != nil { return err }

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		if err := git.InstallPrePushHook(); err != nil {
			return fmt.Errorf("failed to install pre-push hook: %w", err)
		}

		if cfg.EnforceAll {
			if err := git.InstallCommitMsgHook(); err != nil {
				return fmt.Errorf("failed to install commit-msg hook: %w", err)
			}
		}

		fmt.Printf("\n✨ Genie initialized successfully! Config saved to %s\n", config.ConfigFileName)
		return nil
	},
}

func promptOverwrite() bool {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Config file %s already exists. Overwrite?", config.ConfigFileName),
		IsConfirm: true,
	}

	if _, err := prompt.Run(); err != nil {
		return false
	}

	return true;
}

func promptBranch(cfg *config.Config) error {
	promptBranch := promptui.Prompt{
		Label: fmt.Sprintf("What is your protected branch? (default: %s)", cfg.ProtectedBranch),
	}

	protectedBranch, err := promptBranch.Run()
	if err != nil { return err }
	if protectedBranch != "" { cfg.ProtectedBranch = protectedBranch }

	return nil;
}

func promptScope(cfg *config.Config) error {
	promptScope := promptui.Select{
		Label:    "Require scope for all commits?",
		Items:    []string{"No", "Yes"},
		HideHelp: true,
	}

	_, resScope, err := promptScope.Run()
	if err != nil { return err }
	cfg.RequireScope = (resScope == "Yes")

	return nil;
}


func promptEnforce(cfg *config.Config) error {
	promptEnforce := promptui.Select{
		Label:    "Enforce commit standard on all branches or only on protected?",
		Items:    []string{"Pragmatic (protected only)", "Strict (All branches)"},
		HideHelp: true,
	}

	_, resEnforce, err := promptEnforce.Run()
	if err != nil { return err }
	cfg.EnforceAll = (resEnforce == "Strict (All branches)")

	return nil;
}


func promptTypes(cfg *config.Config) error {
	promptTypes := promptui.Prompt{
		Label: "Enter allowed commit types (comma-separated, leave blank for defaults)",
	}

	resTypes, err := promptTypes.Run()
	if err != nil { return err }
	if resTypes != "" {
		parts := strings.Split(resTypes, ",")
		cfg.Types = []string{}
		for _, p := range parts {
			cfg.Types = append(cfg.Types, strings.TrimSpace(p))
		}
	}

	return nil;
}

func promptScopes(cfg *config.Config) error {
	promptScopes := promptui.Prompt{
		Label: "Enter allowed scope names (comma-separated, leave blank for any)",
	}

	resScopes, err := promptScopes.Run()
	if err != nil { return err }
	if resScopes != "" {
		parts := strings.Split(resScopes, ",")
		cfg.AllowedScopes = []string{}
		for _, p := range parts {
			cfg.AllowedScopes = append(cfg.AllowedScopes, strings.TrimSpace(p))
		}
	}

	return nil;
}

const githubWorkflowScriptTemplate = `name: Genie PR Title Check

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

func promptGithubAction(cfg *config.Config) error {
	promptGithubAction := promptui.Prompt{
		Label:    "Create GitHub Actions workflow for PR title enforcement?",
		IsConfirm: true,
		Default: "y",
	}

	if _, err := promptGithubAction.Run(); err != nil {
		return err
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

	script := fmt.Sprintf(githubWorkflowScriptTemplate, cfg.ProtectedBranch)

	if err := os.WriteFile(workflowPath, []byte(script), 0644); err != nil {
		return fmt.Errorf("failed to write workflow file: %w", err)
	}

	fmt.Printf("\n✨ GitHub Action scaffolded successfully at %s\n\n", workflowPath)
	fmt.Println("GitHub Repository Setup Recipe:")
	fmt.Println("1. Go to your repository Settings > Branches.")
	fmt.Println("2. Add a branch protection rule for 'master' (or your main branch).")
	fmt.Println("3. Check 'Require status checks to pass before merging'.")
	fmt.Println("4. Search for 'check-pr-title' and make it required.")
	fmt.Println("5. Check 'Require linear history' or ensure Squash Merging is the only allowed merge strategy.")
	fmt.Println("\nThis ensures that all code landing in your protected branch has a perfectly formatted commit message!")

	return nil;
}
