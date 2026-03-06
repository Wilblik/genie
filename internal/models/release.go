package models

import "time"

// CommitEntry represents a single conventional commit parsed from the history.
type CommitEntry struct {
	Type     string   // e.g., feat, fix, chore
	Scope    string   // e.g., ui, api, db (optional)
	Subject  string   // Short summary
	Body     string   // Detailed explanation (optional)
	Footers  []string // References like "Closes #42" or "BREAKING CHANGE"
	Hash     string   // Git commit hash
	Author   string   // Commit author
	IsBreaking bool   // Derived from '!' in type or footer
}

// ReleaseInfo represents the full data for a single release or changelog range.
type ReleaseInfo struct {
	Tag         string        // The version tag (e.g., v1.0.0)
	Date        time.Time     // Release date
	PreviousTag string        // The tag to compare against
	Modules     []string      // List of all unique scopes detected
	Commits     []CommitEntry // All commits in this release
}
