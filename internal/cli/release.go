package cli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wilblik/genie/internal/config"
	"github.com/wilblik/genie/internal/exporter"
	"github.com/wilblik/genie/internal/git"
	"github.com/wilblik/genie/internal/models"
)

var pushTag bool

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.Flags().BoolVar(&pushTag, "push", false, "Push the newly created tag to the remote (origin)")
}

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Automate a new release (bump version, generate notes, tag)",
	Long:  `Calculates the next semantic version, generates release notes, and creates an annotated Git tag locally.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		latestTag := ""
		tags, _ := git.GetAllTags()
		if len(tags) > 0 {
			latestTag = tags[0]
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("configuration file not found. Please run 'genie init' first")
		}

		releaseInfo, err := getReleaseInfo(latestTag, "HEAD", cfg)
		if err != nil {
			return fmt.Errorf("failed to aggregate commits: %w", err)
		}

		if len(releaseInfo.ChangeTypes) == 0 {
			fmt.Println("No conventional commits found since the last release. Nothing to do.")
			return nil
		}

		var allCommits []models.CommitMessage
		for _, typeGroup := range releaseInfo.ChangeTypes {
			for _, scopes := range typeGroup.Scopes {
				allCommits = append(allCommits, scopes...)
			}
		}

		nextVersion, err := calculateNextVersion(latestTag, allCommits)
		if err != nil {
			return fmt.Errorf("failed to calculate next version: %w", err)
		}

		if nextVersion == latestTag {
			fmt.Println("No features or fixes found to trigger a version bump. Nothing to do.")
			return nil
		}

		fmt.Printf("Genie: Calculated next version as %s\n", nextVersion)

		releaseInfo.Tag = nextVersion
		notes := exporter.GenerateMarkdown(releaseInfo)

		fmt.Printf("Genie: Creating annotated Git tag %s...\n", nextVersion)
		if err := git.CreateTag(nextVersion, notes); err != nil {
			return fmt.Errorf("failed to create tag: %w", err)
		}

		if pushTag {
			fmt.Printf("Genie: Pushing tag %s to origin...\n", nextVersion)
			if err := git.PushTag(nextVersion); err != nil {
				return fmt.Errorf("failed to push tag: %w", err)
			}
			fmt.Printf("✨ Successfully released and pushed %s!\n", nextVersion)
		} else {
			fmt.Printf("✨ Successfully released %s locally!\n", nextVersion)
			fmt.Println("Run 'git push origin --tags' to publish to remote.")
		}

		return nil
	},
}

var semverRegex = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)$`)

func calculateNextVersion(currentTag string, commits []models.CommitMessage) (string, error) {
	if currentTag == "" {
		currentTag = "v0.0.0" // Default for first release
	}

	matches := semverRegex.FindStringSubmatch(strings.TrimSpace(currentTag))
	if matches == nil {
		return "", fmt.Errorf("current tag '%s' is not a valid semantic version (e.g. v1.2.3)", currentTag)
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	const (
		NoBump = iota
		Patch
		Minor
		Major
	)

	impact := NoBump

	for _, c := range commits {
		if c.IsBreaking {
			impact = Major
			break
		}
		if c.ChangeType == "feat" && impact < Minor {
			impact = Minor
		} else if impact < Patch {
			impact = Patch
		}
	}

	switch impact {
	case Major:
		major++; minor = 0; patch = 0
	case Minor:
		minor++; patch = 0
	case Patch:
		patch++
	default:
		return currentTag, nil
	}

	return fmt.Sprintf("v%d.%d.%d", major, minor, patch), nil
}
