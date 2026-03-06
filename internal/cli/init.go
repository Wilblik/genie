package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/wilblik/genie/internal/config"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactive setup for Genie in your repository",
	Long:  `Guides you through setting up Genie, including protected branches, modules, and enforcement levels.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(config.DefaultConfigFileName); err == nil {
			prompt := promptui.Prompt{
				Label:     fmt.Sprintf("Config file %s already exists. Overwrite?", config.DefaultConfigFileName),
				IsConfirm: true,
			}
			if _, err := prompt.Run(); err != nil {
				fmt.Println("Init cancelled.")
				return nil
			}
		}

		cfg := config.NewDefaultConfig()

		promptBranch := promptui.Prompt{
			Label: fmt.Sprintf("What is your protected branch? (default: %s)", cfg.ProtectedBranch),
		}
		resBranch, err := promptBranch.Run()
		if err != nil {
			return err
		}
		if resBranch != "" {
			cfg.ProtectedBranch = resBranch
		}

		promptScope := promptui.Select{
			Label:    "Require module scope for all commits?",
			Items:    []string{"No", "Yes"},
			HideHelp: true,
		}
		_, resScope, err := promptScope.Run()
		if err != nil {
			return err
		}
		cfg.RequireScope = (resScope == "Yes")

		promptEnforce := promptui.Select{
			Label:    "Enforce standard on all branches (Strict) or only on master (Pragmatic)?",
			Items:    []string{"Pragmatic (Master only)", "Strict (All branches)"},
			HideHelp: true,
		}
		_, resEnforce, err := promptEnforce.Run()
		if err != nil {
			return err
		}
		cfg.EnforceAll = (resEnforce == "Strict (All branches)")

		promptModules := promptui.Prompt{
			Label: "Enter allowed module names (comma-separated, leave blank for any)",
		}
		resModules, err := promptModules.Run()
		if err != nil {
			return err
		}
		if resModules != "" {
			parts := strings.Split(resModules, ",")
			cfg.AllowedModules = []string{}
			for _, p := range parts {
				cfg.AllowedModules = append(cfg.AllowedModules, strings.TrimSpace(p))
			}
		}

		promptTypes := promptui.Prompt{
			Label: "Enter allowed commit types (comma-separated, leave blank for defaults)",
		}
		resTypes, err := promptTypes.Run()
		if err != nil {
			return err
		}
		if resTypes != "" {
			parts := strings.Split(resTypes, ",")
			cfg.Types = []string{} // Clear defaults
			for _, p := range parts {
				cfg.Types = append(cfg.Types, strings.TrimSpace(p))
			}
		}

		if err := cfg.Save(config.DefaultConfigFileName); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("\n✨ Genie initialized successfully! Config saved to %s\n", config.DefaultConfigFileName)
		return nil
	},
}
