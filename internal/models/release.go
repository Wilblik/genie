package models

import "time"

// CommitMessage represents a single conventional commit parsed from its raw message string.
type CommitMessage struct {
	ChangeType string            // e.g., feat, fix, chore
	Scope      string            // e.g., ui, api, db (optional)
	Subject    string            // Short summary
	Body       string            // Detailed explanation (optional)
	Footers    map[string]string // Key-value pairs like "Closes": "#42" or "BREAKING CHANGE": "description"
	IsBreaking bool              // Derived from '!' in type or footer
}

// ReleaseInfo represents the full data for a single release or changelog range.
type ReleaseInfo struct {
	Tag    string
	Date   time.Time
	ChangeTypes map[string]ChangeType
}

// Type represents a collection of scopes belonging to a specific commit type.
type ChangeType struct {
	Title  string
	Scopes map[string][]CommitMessage // represents a collection of commits belonging to a specific scope.
}
