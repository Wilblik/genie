package cli

import (
	"testing"

	"github.com/wilblik/genie/internal/models"
)

func TestCalculateNextVersion(t *testing.T) {
	tests := []struct {
		name       string
		currentTag string
		commits    []models.CommitMessage
		want       string
		wantErr    bool
	}{
		{
			name:       "First release (no previous tag)",
			currentTag: "",
			commits: []models.CommitMessage{
				{ChangeType: "feat"},
			},
			want:    "v0.1.0",
			wantErr: false,
		},
		{
			name:       "Patch bump",
			currentTag: "v1.2.3",
			commits: []models.CommitMessage{
				{ChangeType: "fix"},
				{ChangeType: "chore"},
			},
			want:    "v1.2.4",
			wantErr: false,
		},
		{
			name:       "Minor bump (overrides patch)",
			currentTag: "v1.2.3",
			commits: []models.CommitMessage{
				{ChangeType: "fix"},
				{ChangeType: "feat"},
				{ChangeType: "chore"},
			},
			want:    "v1.3.0",
			wantErr: false,
		},
		{
			name:       "Major bump (overrides minor and patch)",
			currentTag: "v1.2.3",
			commits: []models.CommitMessage{
				{ChangeType: "fix"},
				{ChangeType: "feat"},
				{ChangeType: "refactor", IsBreaking: true},
			},
			want:    "v2.0.0",
			wantErr: false,
		},
		{
			name:       "Invalid tag format",
			currentTag: "release-1.0",
			commits: []models.CommitMessage{
				{ChangeType: "fix"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:       "No commits, no bump",
			currentTag: "v1.0.0",
			commits:    []models.CommitMessage{},
			want:       "v1.0.0",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateNextVersion(tt.currentTag, tt.commits)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateNextVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculateNextVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
