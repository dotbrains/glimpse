package git

import (
	"testing"
)

func TestParseNatural_LastNCommits(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"last 3 commits", "HEAD~3"},
		{"last 1 commit", "HEAD~1"},
		{"Last 10 Commits", "HEAD~10"},
	}
	for _, tt := range tests {
		nr := ParseNatural(tt.input)
		if nr == nil {
			t.Errorf("ParseNatural(%q) = nil, want base=%q", tt.input, tt.want)
			continue
		}
		if nr.Base != tt.want {
			t.Errorf("ParseNatural(%q).Base = %q, want %q", tt.input, nr.Base, tt.want)
		}
	}
}

func TestParseNatural_Since(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"since yesterday", "@{1.day.ago}"},
		{"since today", "@{0.day.ago}"},
		{"since 3 days ago", "@{3.days.ago}"},
		{"since 2 weeks ago", "@{2.weeks.ago}"},
	}
	for _, tt := range tests {
		nr := ParseNatural(tt.input)
		if nr == nil {
			t.Errorf("ParseNatural(%q) = nil", tt.input)
			continue
		}
		if nr.Base != tt.want {
			t.Errorf("ParseNatural(%q).Base = %q, want %q", tt.input, nr.Base, tt.want)
		}
	}
}

func TestParseNatural_Focus(t *testing.T) {
	nr := ParseNatural("security issues")
	if nr == nil || nr.Focus == "" {
		t.Fatal("expected focus for 'security issues'")
	}
	if nr.Focus != "security issues" {
		t.Errorf("focus = %q", nr.Focus)
	}

	nr = ParseNatural("performance in src/lib")
	if nr == nil || nr.Focus == "" {
		t.Fatal("expected focus for 'performance in src/lib'")
	}
	if nr.Focus != "performance in src/lib" {
		t.Errorf("focus = %q", nr.Focus)
	}
}

func TestParseNatural_NoMatch(t *testing.T) {
	for _, input := range []string{"main", "HEAD~3", "main..feature", "", "abc1234"} {
		nr := ParseNatural(input)
		if nr != nil {
			t.Errorf("ParseNatural(%q) should be nil, got %+v", input, nr)
		}
	}
}

func TestResolveNaturalOrRefs_FlagsOverride(t *testing.T) {
	base, compare, focus := ResolveNaturalOrRefs([]string{"last 3 commits"}, "develop", "")
	if base != "develop" {
		t.Errorf("flags should override natural language, got base=%q", base)
	}
	if compare != "" {
		t.Errorf("compare = %q", compare)
	}
	if focus != "" {
		t.Errorf("focus = %q", focus)
	}
}

func TestResolveNaturalOrRefs_Natural(t *testing.T) {
	base, compare, _ := ResolveNaturalOrRefs([]string{"last", "5", "commits"}, "", "")
	if base != "HEAD~5" {
		t.Errorf("expected HEAD~5, got %q", base)
	}
	if compare != "" {
		t.Errorf("compare = %q", compare)
	}
}

func TestResolveNaturalOrRefs_FallbackToGitRefs(t *testing.T) {
	base, compare, _ := ResolveNaturalOrRefs([]string{"main..feature"}, "", "")
	if base != "main" || compare != "feature" {
		t.Errorf("expected main/feature, got %q/%q", base, compare)
	}
}
