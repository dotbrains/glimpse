package git

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	lastNCommits = regexp.MustCompile(`(?i)^last\s+(\d+)\s+commits?$`)
	sincePattern = regexp.MustCompile(`(?i)^since\s+(yesterday|today|\d+\s+(?:days?|weeks?|months?)\s+ago)$`)
)

// NaturalRef represents a parsed natural language reference.
type NaturalRef struct {
	Base    string
	Compare string
	Focus   string // extracted focus area if present
}

// ParseNatural tries to interpret a natural language string as git refs.
// Returns nil if the input doesn't match any known patterns.
func ParseNatural(input string) *NaturalRef {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	// "last N commits"
	if m := lastNCommits.FindStringSubmatch(input); m != nil {
		n, _ := strconv.Atoi(m[1])
		if n > 0 {
			return &NaturalRef{Base: "HEAD~" + strconv.Itoa(n)}
		}
	}

	// "since yesterday" / "since 3 days ago"
	if m := sincePattern.FindStringSubmatch(input); m != nil {
		period := m[1]
		switch {
		case strings.EqualFold(period, "yesterday"):
			return &NaturalRef{Base: "@{1.day.ago}"}
		case strings.EqualFold(period, "today"):
			return &NaturalRef{Base: "@{0.day.ago}"}
		default:
			// "3 days ago" -> "@{3.days.ago}"
			period = strings.ReplaceAll(period, " ", ".")
			return &NaturalRef{Base: "@{" + period + "}"}
		}
	}

	// Check for focus areas embedded in natural language:
	// "security issues" / "performance in src/lib" / "testing"
	focus, path := extractFocus(input)
	if focus != "" {
		ref := &NaturalRef{Focus: focus}
		if path != "" {
			ref.Focus = focus + " in " + path
		}
		return ref
	}

	return nil
}

// extractFocus looks for known focus keywords, optionally followed by "in <path>".
func extractFocus(input string) (focus, path string) {
	lower := strings.ToLower(input)

	focusKeywords := []string{
		"security issues", "security",
		"performance issues", "performance",
		"testing", "tests",
		"error handling", "errors",
		"memory", "concurrency", "race conditions",
	}

	for _, kw := range focusKeywords {
		if strings.Contains(lower, kw) {
			focus = kw
			// Check for "in <path>" suffix.
			idx := strings.Index(lower, " in ")
			if idx > 0 && idx > strings.Index(lower, kw) {
				path = strings.TrimSpace(input[idx+4:])
			}
			return
		}
	}
	return "", ""
}

// ResolveNaturalOrRefs tries natural language parsing first, then falls back to git ref parsing.
func ResolveNaturalOrRefs(args []string, flagBase, flagCompare string) (base, compare, focus string) {
	// Flags always take priority.
	if flagBase != "" {
		return flagBase, flagCompare, ""
	}

	// Try natural language on the full joined input.
	if len(args) > 0 {
		joined := strings.Join(args, " ")
		if nr := ParseNatural(joined); nr != nil {
			return nr.Base, nr.Compare, nr.Focus
		}
	}

	// Fall back to standard git ref parsing.
	b, c := ResolveRefs(args, flagBase, flagCompare)
	return b, c, ""
}
