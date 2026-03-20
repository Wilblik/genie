package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const PrePushHookScript = `#!/bin/sh
# Genie: Automated Release Notes & Commit Enforcement
# This hook calls Genie to validate your commits against the standard.

# On some systems, the PATH might not include where Genie is installed.
# We try to call 'genie' directly.
genie check-push "$@" <&0
`

const CommitMsgHookScript = `#!/bin/sh
# Genie: Automated Release Notes & Commit Enforcement
# This hook calls Genie to validate your commit message.

msg=$(cat "$1")
genie check-msg "$msg"
`

func InstallPrePushHook() error {
	return installHook("pre-push", PrePushHookScript)
}

func InstallCommitMsgHook() error {
	return installHook("commit-msg", CommitMsgHookScript)
}

/* GetCommitMessages returns the commit messages for a given range.
 * from/to can be tags, hashes, branches, or dates (YYYY-MM-DD).
**/
func GetCommitMessages(from, to string) ([]string, error) {
	args := []string{"log", "--format=%B%x00"}

	if from == "TAIL" || from == "0000000000000000000000000000000000000000" {
		from = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"
	}

	isFromDate := isDate(from)
	if from != "" {
		if isFromDate {
			args = append(args, "--since="+formatDate(from))
		} else {
			// TODO Is it needed?
			//if err := exec.Command("git", "cat-file", "-e", from).Run(); err != nil {
			//	return nil, fmt.Errorf("starting point '%s' not found in git history", from)
			//}
			args = append(args, from)
		}
	}


	if to == "" { to = "HEAD" }
	if isDate(to) {
		args = append(args, "--until="+formatDate(to))
	} else if from != "" && !isFromDate {
		fromTag := args[len(args)-1]
		args[len(args)-1] = fromTag+".."+to
	} else {
		args = append(args, to)
	}

	fmt.Println(args)
	return runGitLog(args)
}

func GetAllTags() ([]string, error) {
	// --sort=-v:refname sorts tags by version in descending order (v1.10 > v1.2)
	cmd := exec.Command("git", "tag", "--sort=-v:refname")
	out, err := cmd.Output()
	if err != nil { return nil, fmt.Errorf("failed to get tags: %w", err) }

	var result []string
	for t := range strings.SplitSeq(string(out), "\n") {
		t = strings.TrimSpace(t)
		if t != "" {
			result = append(result, t)
		}
	}
	return result, nil
}

func CreateTag(tag, message string) error {
	cmd := exec.Command("git", "tag", "-a", tag, "-m", message)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git tag failed: %v\nOutput: %s", err, string(out))
	}
	return nil
}

func PushTag(tag string) error {
	cmd := exec.Command("git", "push", "origin", tag)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git push tag failed: %v\nOutput: %s", err, string(out))
	}
	return nil
}

func isDate(s string) bool {
	match, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}( \d{2}:\d{2}:\d{2})?$`, s) // YYYY-MM-DD HH:MM:SS
	return match
}

func runGitLog(args []string) ([]string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git log failed: %w", err)
	}

	raw := string(output)
	if raw == "" {
		return nil, nil
	}

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

func installHook(name, script string) error {
	hooksDir, err := getHooksDir()
	if err != nil { return err; }

	hookPath := filepath.Join(hooksDir, name)
	if err := os.WriteFile(hookPath, []byte(script), 0755); err != nil {
		return fmt.Errorf("failed to write hook file %s: %w", name, err)
	}

	return nil
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

func formatDate(date string) string {
	newDate := date
	t, err := time.Parse("2006-01-02 15:04:05", newDate)
	if err == nil {
		newDate = t.Format("2006-01-02 15:04:05")
	} else {
		t, err = time.Parse("2006-01-02", newDate)
		if err == nil {
			newDate = t.Format("2006-01-02 15:04:05")
		}
	}

	return newDate
}
