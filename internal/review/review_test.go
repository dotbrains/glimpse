package review

import (
	"testing"
)

func TestParseReviewOutput_ValidComments(t *testing.T) {
	raw := `main.go:42 [must-fix] Nil pointer dereference — check err before using resp.Body
auth.go:18 [suggestion] Consider using sync.Once for the token refresh to avoid races
config.go:5 [nit] Exported function missing godoc comment
handler.go:99 [question] Is this timeout intentional or should it be configurable?`

	comments := ParseReviewOutput(raw)
	if len(comments) != 4 {
		t.Fatalf("expected 4 comments, got %d", len(comments))
	}

	c := comments[0]
	if c.File != "main.go" {
		t.Errorf("file = %q, want main.go", c.File)
	}
	if c.Line != 42 {
		t.Errorf("line = %d, want 42", c.Line)
	}
	if c.Severity != "must-fix" {
		t.Errorf("severity = %q, want must-fix", c.Severity)
	}
	if c.Author != "ai" {
		t.Errorf("author = %q, want ai", c.Author)
	}

	if comments[1].Severity != "suggestion" {
		t.Errorf("comment[1] severity = %q", comments[1].Severity)
	}
	if comments[2].Severity != "nit" {
		t.Errorf("comment[2] severity = %q", comments[2].Severity)
	}
	if comments[3].Severity != "question" {
		t.Errorf("comment[3] severity = %q", comments[3].Severity)
	}
}

func TestParseReviewOutput_Empty(t *testing.T) {
	comments := ParseReviewOutput("")
	if len(comments) != 0 {
		t.Errorf("expected 0 comments for empty input, got %d", len(comments))
	}
}

func TestParseReviewOutput_MixedLines(t *testing.T) {
	raw := `Here is my review of the code:

main.go:10 [must-fix] Missing error check
Some non-matching preamble text.
lib.go:20 [nit] Rename this variable

That's all I found.`

	comments := ParseReviewOutput(raw)
	if len(comments) != 2 {
		t.Fatalf("expected 2 comments (ignoring non-matching lines), got %d", len(comments))
	}
	if comments[0].File != "main.go" || comments[0].Line != 10 {
		t.Errorf("first = %s:%d", comments[0].File, comments[0].Line)
	}
	if comments[1].File != "lib.go" || comments[1].Line != 20 {
		t.Errorf("second = %s:%d", comments[1].File, comments[1].Line)
	}
}

func TestParseReviewOutput_PathsWithSlashes(t *testing.T) {
	raw := `src/handlers/auth.go:55 [suggestion] Extract this into a helper function`

	comments := ParseReviewOutput(raw)
	if len(comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(comments))
	}
	if comments[0].File != "src/handlers/auth.go" {
		t.Errorf("file = %q", comments[0].File)
	}
}

func TestParseReviewOutput_InvalidSeverity(t *testing.T) {
	// Only the 4 valid severities should match.
	raw := `main.go:1 [critical] This won't match the pattern`
	comments := ParseReviewOutput(raw)
	if len(comments) != 0 {
		t.Errorf("expected 0 comments for invalid severity, got %d", len(comments))
	}
}
