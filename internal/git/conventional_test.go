package git

import (
	"reflect"
	"testing"

	"github.com/wilblik/genie/internal/config"
)

func TestParseCommit(t *testing.T) {
	tests := []struct {
		name        string
		raw         string
		wantType    string
		wantSubj    string
		wantScope   string
		wantBody    string
		isBreaking  bool
		wantFooters map[string]string
	}{
		{
			name:     "Simple feat",
			raw:      "feat: add new feature",
			wantType: "feat",
			wantSubj: "add new feature",
		},
		{
			name:      "Feat with scope",
			raw:       "feat(ui): add button",
			wantType:  "feat",
			wantSubj:  "add button",
			wantScope: "ui",
		},
		{
			name:       "Breaking change with !",
			raw:        "feat(api)!: remove endpoint",
			wantType:   "feat",
			wantSubj:   "remove endpoint",
			wantScope:  "api",
			isBreaking: true,
		},
		{
			name:        "Breaking change in footer",
			raw:         "feat: add feature\n\nBREAKING CHANGE: this is breaking",
			wantType:    "feat",
			wantSubj:    "add feature",
			isBreaking:  true,
			wantFooters: map[string]string{"BREAKING CHANGE": "this is breaking"},
		},
		{
			name:        "Custom footer",
			raw:         "fix(core): resolve race condition\n\nJira-Ticket: PROJ-123\nSigned-off-by: Test User",
			wantType:    "fix",
			wantSubj:    "resolve race condition",
			wantScope:   "core",
			wantFooters: map[string]string{"Jira-Ticket": "PROJ-123", "Signed-off-by": "Test User"},
		},
		{
			name:        "Body and footers",
			raw:         "feat: some feat\n\nThis is the body content.\n\nRefs: #42",
			wantType:    "feat",
			wantSubj:    "some feat",
			wantBody:    "This is the body content.",
			wantFooters: map[string]string{"Refs": "#42"},
		},
		{
			name:        "Footer with # separator",
			raw:         "fix: some fix\n\nFixes #123",
			wantType:    "fix",
			wantSubj:    "some fix",
			wantFooters: map[string]string{"Fixes": "123"},
		},
		{
			name:       "BREAKING CHANGE in body (not footer)",
			raw:        "feat: subject\n\nBREAKING CHANGE: this is in the body\nThis makes it a body, not a footer.",
			wantType:   "feat",
			wantSubj:   "subject",
			wantBody:   "BREAKING CHANGE: this is in the body\nThis makes it a body, not a footer.",
			isBreaking: false,
		},
	}

	cfg := config.NewDefaultConfig()
	cfg.Types = []string{"feat", "fix"}
	cfg.AllowedScopes = []string{"ui", "api", "core"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCommitMessage(tt.raw, cfg)
			if err != nil {
				t.Errorf("ParsecommitMessage() error = %v", err)
				return
			}
			if got.ChangeType != tt.wantType {
				t.Errorf("Type = %v, want %v", got.ChangeType, tt.wantType)
			}
			if got.Subject != tt.wantSubj {
				t.Errorf("Subject = %v, want %v", got.Subject, tt.wantSubj)
			}
			if got.Scope != tt.wantScope {
				t.Errorf("Scope = %v, want %v", got.Scope, tt.wantScope)
			}
			if got.Body != tt.wantBody {
				t.Errorf("Body = %v, want %v", got.Body, tt.wantBody)
			}
			if got.IsBreaking != tt.isBreaking {
				t.Errorf("IsBreaking = %v, want %v", got.IsBreaking, tt.isBreaking)
			}
			if len(got.Footers) == 0 && len(tt.wantFooters) == 0 {
				return
			}
			if !reflect.DeepEqual(got.Footers, tt.wantFooters) {
				t.Errorf("Footers = %v, want %v", got.Footers, tt.wantFooters)
			}
		})
	}
}
