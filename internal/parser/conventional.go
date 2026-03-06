package parser

import (
	"regexp"
	"strings"

	"github.com/wilblik/genie/internal/models"
)

var (
	// Regex for the first line: type(scope)!: subject
	headerRegex = regexp.MustCompile(`^([a-z]+)(?:\(([^)]+)\))?(!)?:\s+(.+)$`)
	// Regex for footers: Token: Value or Token #Value
	footerRegex = regexp.MustCompile(`(?m)^([a-zA-Z0-9-]+|BREAKING CHANGE)(: | #).+$`)
)

// ParseCommit parses a raw git commit message string into a CommitEntry model.
// It uses string slicing to create "views" into the original message, minimizing allocations.
func ParseCommit(raw string) (*models.CommitEntry, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	// 1. Identify the Header (everything up to the first newline)
	headerEnd := strings.IndexByte(raw, '\n')
	header := raw
	if headerEnd != -1 {
		header = raw[:headerEnd]
	}

	header_parts := headerRegex.FindStringSubmatch(header)
	if header_parts == nil {
		return nil, nil
	}

	entry := &models.CommitEntry{
		Type:       header_parts[1],
		Scope:      header_parts[2],
		IsBreaking: header_parts[3] == "!",
		Subject:    header_parts[4],
	}

	if headerEnd == -1 {
		return entry, nil
	}

	// 2. Identify the Footer block
	// Footers are at the very end of the message, separated from the body by a blank line.
	footerStart := len(raw)
	lastBlankLine := strings.LastIndex(raw, "\n\n")

	if lastBlankLine != -1 {
		potentialFooterBlock := raw[lastBlankLine+2:]
		lines := strings.Split(potentialFooterBlock, "\n")
		isAllFooters := true
		for _, line := range lines {
			if line != "" && !footerRegex.MatchString(line) {
				isAllFooters = false
				break
			}
		}

		if isAllFooters {
			footerStart = lastBlankLine + 2
			entry.Footers = lines
			for _, f := range lines {
				if strings.HasPrefix(strings.ToUpper(f), "BREAKING CHANGE:") {
					entry.IsBreaking = true
				}
			}
		}
	}

	// 3. Identify the Body
	bodyStart := headerEnd + 1
	if bodyStart < len(raw) && raw[bodyStart] == '\n' {
		bodyStart++
	}

	if bodyStart < footerStart {
		entry.Body = strings.TrimSpace(raw[bodyStart:footerStart])
	}

	return entry, nil
}
