package exporter

import (
	"fmt"
	"strings"

	"github.com/wilblik/genie/internal/models"
)

func GenerateMarkdown(release *models.ReleaseInfo) string {
	var sb strings.Builder
	var breakingSb strings.Builder

	fmt.Fprintf(&sb, "# Release %s (%s)\n\n", release.Tag, release.Date.Format("2006-01-02"))

	for _, group := range release.Groups {
		fmt.Fprintf(&sb, "## %s\n", group.Title)

		for scope, commits := range group.Scopes {
			fmt.Fprintf(&sb, "### %s\n", scope)
			for _, commit := range commits {
				breakingBadge := ""
				if commit.IsBreaking {
					breakingBadge = " ⚠️ **BREAKING**"
					if breakingChangeDesc, ok := commit.Footers["BREAKING CHANGE"]; ok {
						fmt.Fprintf(&breakingSb, "- **%s**: %s\n", commit.Subject, breakingChangeDesc)
					}
				}

				fmt.Fprintf(&sb, "- %s%s\n", commit.Subject, breakingBadge)
			}
		}
		sb.WriteString("\n")
	}

	if breakingSb.Len() > 0 {
		sb.WriteString("## ⚠️ BREAKING CHANGES\n")
		sb.WriteString(breakingSb.String() + "\n")
	}

	return sb.String()
}
