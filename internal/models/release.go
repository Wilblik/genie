package models

import "time"

// CommitMessage represents a single conventional commit parsed from its raw message string.
type CommitMessage struct {
	Type       string            // e.g., feat, fix, chore
	Scope      string            // e.g., ui, api, db (optional)
	Subject    string            // Short summary
	Body       string            // Detailed explanation (optional)
	Footers    map[string]string // Key-value pairs like "Closes": "#42" or "BREAKING CHANGE": "description"
	IsBreaking bool              // Derived from '!' in type or footer
}

// ReleaseInfo represents the full data for a single release or changelog range.
type ReleaseInfo struct {
	Tag         string          // The version tag (e.g., v1.0.0)
	Date        time.Time       // Release date
	Modules     []string        // List of all unique scopes detected
	Commits     []CommitMessage // All commits in this release
}
