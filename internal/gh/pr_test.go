package gh

import (
	"testing"
)

func TestIsPRURL_Valid(t *testing.T) {
	urls := []string{
		"https://github.com/owner/repo/pull/123",
		"https://github.com/dotbrains/glimpse/pull/1",
		"http://github.com/org/project/pull/99999",
		"https://github.com/a/b/pull/1",
	}
	for _, u := range urls {
		if !IsPRURL(u) {
			t.Errorf("IsPRURL(%q) = false, want true", u)
		}
	}
}

func TestIsPRURL_Invalid(t *testing.T) {
	urls := []string{
		"https://github.com/owner/repo",
		"https://github.com/owner/repo/issues/123",
		"https://gitlab.com/owner/repo/pull/1",
		"main..feature",
		"HEAD~3",
		"",
		"not a url",
	}
	for _, u := range urls {
		if IsPRURL(u) {
			t.Errorf("IsPRURL(%q) = true, want false", u)
		}
	}
}

func TestParsePRURL_Valid(t *testing.T) {
	owner, repo, number, err := ParsePRURL("https://github.com/dotbrains/glimpse/pull/42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if owner != "dotbrains" {
		t.Errorf("owner = %q, want dotbrains", owner)
	}
	if repo != "glimpse" {
		t.Errorf("repo = %q, want glimpse", repo)
	}
	if number != "42" {
		t.Errorf("number = %q, want 42", number)
	}
}

func TestParsePRURL_Invalid(t *testing.T) {
	_, _, _, err := ParsePRURL("https://github.com/owner/repo")
	if err == nil {
		t.Fatal("expected error for non-PR URL")
	}
}

func TestGHInstalled(t *testing.T) {
	// Just ensure it doesn't panic.
	_ = GHInstalled()
}
