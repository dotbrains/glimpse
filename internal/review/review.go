package review

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/dotbrains/glimpse/internal/comments"
	gexec "github.com/dotbrains/glimpse/internal/exec"
)

var commentPattern = regexp.MustCompile(`^([^:]+):(\d+)\s+\[(must-fix|suggestion|nit|question)\]\s+(.+)$`)

// Reviewer runs AI code review on a diff.
type Reviewer struct {
	Executor gexec.CommandExecutor
}

// NewReviewer creates a Reviewer.
func NewReviewer() *Reviewer {
	return &Reviewer{Executor: gexec.NewRealExecutor()}
}

const reviewPrompt = `You are a senior code reviewer. Review the following git diff and output inline comments.

Rules:
- One comment per line, format: file:line [severity] comment text
- Severity must be one of: [must-fix], [suggestion], [nit], [question]
- Be specific and actionable. Reference exact variable/function names.
- No preamble, no summary, just the comment lines.
- Skip praise. Only output issues.
- Write like a senior engineer — direct, specific, no AI-speak.

%sDiff:
%s`

// Run executes an AI review and returns parsed comments.
func (r *Reviewer) Run(ctx context.Context, rawDiff string, focus string) ([]comments.Comment, error) {
	focusLine := ""
	if focus != "" {
		focusLine = fmt.Sprintf("Focus on: %s\n\n", focus)
	}

	prompt := fmt.Sprintf(reviewPrompt, focusLine, rawDiff)

	out, err := r.Executor.RunWithStdin(ctx, prompt, "claude", "-p", "-")
	if err != nil {
		return nil, fmt.Errorf("AI review failed: %w", err)
	}

	return ParseReviewOutput(out), nil
}

// ParseReviewOutput parses structured review output into Comment structs.
func ParseReviewOutput(raw string) []comments.Comment {
	var result []comments.Comment

	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		m := commentPattern.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		lineNum, _ := strconv.Atoi(m[2])
		result = append(result, comments.Comment{
			File:     m[1],
			Line:     lineNum,
			Side:     "new",
			Severity: m[3],
			Body:     m[4],
			Author:   "ai",
		})
	}
	return result
}
