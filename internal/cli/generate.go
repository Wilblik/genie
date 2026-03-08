package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/wilblik/genie/internal/config"
	"github.com/wilblik/genie/internal/exporter"
	"github.com/wilblik/genie/internal/git"
	"github.com/wilblik/genie/internal/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	from   string
	to     string
	format string
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVar(&from, "from", "", "Starting point (tag, hash, or YYYY-MM-DD)")
	generateCmd.Flags().StringVar(&to, "to", "HEAD", "Ending point (tag, hash, or branch)")
	generateCmd.Flags().StringVar(&format, "format", "markdown", "Output format (markdown)")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate release notes from Git history",
	Long:  `Analyzes Git commits and generates a formatted release report.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if from == "" {
			tags, _ := git.GetAllTags()
			if len(tags) > 0 {
				if to != tags[0] {
					from = tags[0]
				} else if len(tags) > 1 {
					from = tags[1]
				}
			}
		}

		if from == "" {
			fmt.Printf("Genie: Generating release notes for entire history up to %s...\n", to)
		} else {
			fmt.Printf("Genie: Generating release notes from %s to %s...\n", from, to)
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("configuration file reading failed. Please run 'genie init' first if %s is not present", config.ConfigFileName)
		}

		release, err := getReleaseInfo(from, to, cfg)
		if err != nil {
			return fmt.Errorf("failed to aggregate commits: %w", err)
		}

		if format != "markdown" {
			return fmt.Errorf("format '%s' is not supported", format)
		}

		output := exporter.GenerateMarkdown(release)
		fmt.Println("\n" + output)

		return nil
	},
}

func getReleaseInfo(from, to string, cfg *config.Config) (*models.ReleaseInfo, error) {
	commitMessages, err := git.GetCommitMessages(from, to)
	if err != nil { return nil, err; }

	releaseInfo := &models.ReleaseInfo{
		Tag:  to,
		Date: time.Now(),
		Groups: make(map[string]models.TypeGroup),
	}

	for _, msg := range commitMessages {
		// Skip non conventional commits
		parsedCommitMsg, err := git.ParseCommitMessage(msg)
		if err != nil || parsedCommitMsg == nil { continue }
		if err := git.ValidateCommitMessage(cfg, parsedCommitMsg); err != nil { continue }

		scope := parsedCommitMsg.Scope
		if scope == "" { scope = "general" }

		if _, ok := releaseInfo.Groups[parsedCommitMsg.Type]; !ok {
			releaseInfo.Groups[parsedCommitMsg.Type] = models.TypeGroup{
				Title: getTitleForType(parsedCommitMsg.Type),
				Scopes: make(map[string][]models.CommitMessage),
			}
		}

		releaseInfo.Groups[parsedCommitMsg.Type].Scopes[scope] = append(
			releaseInfo.Groups[parsedCommitMsg.Type].Scopes[scope],
			*parsedCommitMsg,
		)
	}

	return releaseInfo, nil
}

func getTitleForType(t string) string {
	titles := map[string]string{
		"feat":     "Features",
		"fix":      "Bug Fixes",
		"perf":     "Performance",
		"refactor": "Refactoring",
		"docs":     "Documentation",
		"chore":    "Miscellaneous",
		"test":     "Tests",
		"style":    "Style",
		"revert":   "Reverts",
	}

	if title, ok := titles[t]; ok {
		return title
	}

	caser := cases.Title(language.English)
	return caser.String(t)
}
