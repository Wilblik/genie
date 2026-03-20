package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/wilblik/genie/internal/config"
	"github.com/wilblik/genie/internal/git"
	"github.com/wilblik/genie/internal/forge"
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

func promptGithubAction(cfg *config.Config) error {
	promptPR := promptui.Select{
		Label:     "Create GitHub Action to enforce Pull Requests Titles?",
		Items:    []string{"No", "Yes"},
		HideHelp: true,
	}

	_, promptPrResult, promptPrErr := promptPR.Run()
	if promptPrErr != nil { return promptPrErr }
	createPr := promptPrResult == "Yes"

	promptRelease := promptui.Select{
		Label:     "Create GitHub Action for Automated Releases (on demand)?",
		Items:    []string{"No", "Yes"},
		HideHelp: true,
	}

	_, promptReleaseResult, promptReleaseErr := promptRelease.Run()
	if promptReleaseErr != nil { return promptReleaseErr }
	createRelease := promptReleaseResult == "Yes"

	if !createPr && !createRelease { return nil }

	workflowDir, err := forge.CreateGithubWorkflowDir()
	if err != nil { return err }

	if createPr {
		err := forge.CreateGithubPullRequestWorkflow(workflowDir, cfg.ProtectedBranch)
		if err != nil { return err }
	}

	if createRelease {
		err := forge.CreateGithubReleaseWorkflow(workflowDir)
		if err != nil { return err }
	}

	return nil
}
