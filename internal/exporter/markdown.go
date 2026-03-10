package exporter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/wilblik/genie/internal/models"
)

// Any custom types will be sorted alphabetically after these.
var preferredOrder = map[string]int{
	"feat":     1,
	"fix":      2,
	"perf":     3,
	"refactor": 4,
	"docs":     5,
	"style":    6,
	"test":     7,
	"chore":    8,
	"revert":   9,
}

func GenerateMarkdown(release *models.ReleaseInfo) string {
	var sb strings.Builder
	var breakingSb strings.Builder

	fmt.Fprintf(&sb, "# Release %s (%s)\n\n", release.Tag, release.Date.Format("2006-01-02"))

	types := sortReleaseInfoByType(release)

	for _, t := range types {
		typeGroup := release.ChangeTypes[t]
		fmt.Fprintf(&sb, "## %s\n", typeGroup.Title)

		scopeNames := sortScopes(&typeGroup)
		for _, scopeName := range scopeNames {
			printCommitsByScope(&sb, &breakingSb, scopeName, &typeGroup)
		}

		sb.WriteString("\n")
	}

	if breakingSb.Len() > 0 {
		sb.WriteString("## ⚠️ BREAKING CHANGES\n")
		sb.WriteString(breakingSb.String() + "\n")
	}

	sb.WriteString("\n")
	return strings.TrimSpace(sb.String())
}

func sortReleaseInfoByType(release *models.ReleaseInfo) []string {
	var types []string
	for t := range release.ChangeTypes {
		types = append(types, t)
	}

	sort.Slice(types, func(i, j int) bool {
		t1, t2 := types[i], types[j]
		p1, ok1 := preferredOrder[t1]
		p2, ok2 := preferredOrder[t2]

		// Standard types come first
		if ok1 && ok2 { return p1 < p2 }
		if ok1 { return true }
		if ok2 { return false }
		return t1 < t2 // Custom types sorted alphabetically
	})

	return types
}


func sortScopes(typeGroup *models.ChangeType) []string {
	var scopes []string
	for s := range typeGroup.Scopes {
		scopes = append(scopes, s)
	}
	sort.Strings(scopes)
	return scopes
}

func printCommitsByScope(sb *strings.Builder, breakingSb *strings.Builder, scopeName string, typeGroup *models.ChangeType) {
	if scopeName != "" {
		fmt.Fprintf(sb, "### %s\n", scopeName)
	}

	commits := typeGroup.Scopes[scopeName]
	for _, commit := range commits {
		breakingBadge := ""
		if commit.IsBreaking {
			breakingBadge = " ⚠️ **BREAKING**"
			if breakingChangeDesc, ok := commit.Footers["BREAKING CHANGE"]; ok {
				fmt.Fprintf(breakingSb, "- **%s**: %s\n", commit.Subject, breakingChangeDesc)
			}
		}
		fmt.Fprintf(sb, "- %s%s\n", commit.Subject, breakingBadge)
	}

	if scopeName != "" {
		sb.WriteString("\n")
	}
}
