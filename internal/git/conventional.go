package git

import (
	"fmt"
	"regexp"
	"strings"
	"slices"

	"github.com/wilblik/genie/internal/models"
	"github.com/wilblik/genie/internal/config"
)

var (
	// Regex for the first line: type(scope)!: subject
	headerRegex = regexp.MustCompile(`^([a-z]+)(?:\(([^)]+)\))?(!)?:\s+(.+)$`)
	// Regex for footers: (Token)(: | #)Value
	footerRegex = regexp.MustCompile(`^([a-zA-Z0-9-]+|BREAKING CHANGE)(: | #)`)
)

func ParseCommitMessage(raw string, cfg *config.Config) (*models.CommitMessage, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("commit message is empty")
	}

	commitMsg := &models.CommitMessage{}

	success, headerEnd := parseHeader(raw, commitMsg)
	if !success {
		return nil, fmt.Errorf("message does not follow Conventional Commits standard\nExample: feat(ui): add new button")
	}

	if !slices.Contains(cfg.Types, commitMsg.ChangeType) {
		return nil, fmt.Errorf("type '%s' is not in the allowed types list: %v", commitMsg.ChangeType, cfg.Types)
	}

	if cfg.RequireScope && commitMsg.Scope == "" {
		return nil, fmt.Errorf("scope is required but was not provided")
	}

	if len(cfg.AllowedScopes) > 0 && commitMsg.Scope != "" {
		if !slices.Contains(cfg.AllowedScopes, commitMsg.Scope) {
			return nil, fmt.Errorf("scope '%s' is not in the allowed scopes list: %v", commitMsg.Scope, cfg.AllowedScopes)
		}
	}

	if headerEnd == -1 {
		return commitMsg, nil
	}

	footerStart := parseFooter(raw, commitMsg)
	parseBody(raw, commitMsg, headerEnd, footerStart)

	return commitMsg, nil
}

func parseHeader(raw string, commitMsg *models.CommitMessage) (bool, int) {
	headerEnd := strings.IndexByte(raw, '\n')
	header := raw
	if headerEnd != -1 {
		header = raw[:headerEnd]
	}

	header_parts := headerRegex.FindStringSubmatch(header)
	if header_parts == nil {
		return false, headerEnd
	}

	commitMsg.ChangeType = header_parts[1]
	commitMsg.Scope = header_parts[2]
	commitMsg.IsBreaking = header_parts[3] == "!"
	commitMsg.Subject = header_parts[4]

	return true, headerEnd
}

func parseFooter(raw string, commitMsg *models.CommitMessage) int {
	footerStart := len(raw)
	lastBlankLine := strings.LastIndex(raw, "\n\n")

	if lastBlankLine != -1 {
		potentialFooterBlock := raw[lastBlankLine+2:]
		lines := strings.Split(potentialFooterBlock, "\n")
		tempFooters := make(map[string]string)
		isValidBlock := true

		for _, line := range lines {
			if line == "" {
				continue
			}
			parts := footerRegex.FindStringSubmatch(line)
			if parts == nil {
				isValidBlock = false
				break
			}

			key := parts[1]
			sep := parts[2]
			value := line[len(key)+len(sep):]
			tempFooters[key] = value
		}

		if isValidBlock && len(tempFooters) > 0 {
			footerStart = lastBlankLine + 2
			commitMsg.Footers = tempFooters
			if _, ok := tempFooters["BREAKING CHANGE"]; ok {
				commitMsg.IsBreaking = true
			}
		}
	}

	return footerStart
}

func parseBody(raw string, commitMsg *models.CommitMessage, headerEnd int, footerStart int) {
	bodyStart := headerEnd + 1
	if bodyStart < len(raw) && raw[bodyStart] == '\n' {
		bodyStart++
	}

	if bodyStart < footerStart {
		commitMsg.Body = strings.TrimSpace(raw[bodyStart:footerStart])
	}
}
