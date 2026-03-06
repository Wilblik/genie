package parser

import (
	"reflect"
	"testing"
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
		wantFooters []string
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
			wantFooters: []string{"BREAKING CHANGE: this is breaking"},
		},
		{
			name:        "Custom footer",
			raw:         "fix(core): resolve race condition\n\nJira-Ticket: PROJ-123\nSigned-off-by: Test User",
			wantType:    "fix",
			wantSubj:    "resolve race condition",
			wantScope:   "core",
			wantFooters: []string{"Jira-Ticket: PROJ-123", "Signed-off-by: Test User"},
		},
		{
			name:        "Body and footers",
			raw:         "feat: some feat\n\nThis is the body content.\n\nRefs: #42",
			wantType:    "feat",
			wantSubj:    "some feat",
			wantBody:    "This is the body content.",
			wantFooters: []string{"Refs: #42"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCommit(tt.raw)
			if err != nil {
				t.Fatalf("ParseCommit() error = %v", err)
			}
			if got == nil {
				t.Fatal("ParseCommit() returned nil")
			}
			if got.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", got.Type, tt.wantType)
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
