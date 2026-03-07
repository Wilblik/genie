package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const PrePushHookScript = `#!/bin/sh
# Genie: Automated Release Notes & Commit Enforcement
# This hook calls Genie to validate your commits against the standard.

# On some systems, the PATH might not include where Genie is installed.
# We try to call 'genie' directly.
genie check-push "$@" <&0
`

func InstallPrePushHook() error {
	hooksDir, err := getHooksDir()
	if err != nil { return err; }

	hookPath := filepath.Join(hooksDir, "pre-push")
	if err := os.WriteFile(hookPath, []byte(PrePushHookScript), 0755); err != nil {
		return fmt.Errorf("failed to write hook file: %w", err)
	}

	return nil
}

func GetCommitMessages(from, to string) ([]string, error) {
	rangeSpec := to
	if from != "" && from != "0000000000000000000000000000000000000000" {
		rangeSpec = fmt.Sprintf("%s..%s", from, to)
	}

	cmd := exec.Command("git", "log", rangeSpec, "--format=%B%x00")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git log failed: %w", err)
	}

	raw := string(output)
	if raw == "" {
		return nil, nil
	}

	// Split by null byte and trim the last empty element
	messages := strings.Split(raw, "\x00")
	var result []string
	for _, m := range messages {
		m = strings.TrimSpace(m)
		if m != "" {
			result = append(result, m)
		}
	}

	return result, nil
}

func getHooksDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil { return "", err; }

	curr := wd
	for {
		gitDir := filepath.Join(curr, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return filepath.Join(gitDir, "hooks"), nil
		}

		parent := filepath.Dir(curr)
		if parent == curr { break; }
		curr = parent
	}

	return "", fmt.Errorf("not a git repository")
}
