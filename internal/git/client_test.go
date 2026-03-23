package git

import (
	"testing"
)

func TestResolveRefs_Empty(t *testing.T) {
	base, compare := ResolveRefs(nil, "", "")
	if base != "" || compare != "" {
		t.Errorf("expected empty refs, got base=%q compare=%q", base, compare)
	}
}

func TestResolveRefs_SingleRef(t *testing.T) {
	base, compare := ResolveRefs([]string{"main"}, "", "")
	if base != "main" || compare != "" {
		t.Errorf("expected base=main compare='', got base=%q compare=%q", base, compare)
	}
}

func TestResolveRefs_TwoArgs(t *testing.T) {
	base, compare := ResolveRefs([]string{"main", "feature"}, "", "")
	if base != "main" || compare != "feature" {
		t.Errorf("expected base=main compare=feature, got base=%q compare=%q", base, compare)
	}
}

func TestResolveRefs_RangeSyntax(t *testing.T) {
	base, compare := ResolveRefs([]string{"main..feature"}, "", "")
	if base != "main" || compare != "feature" {
		t.Errorf("expected base=main compare=feature, got base=%q compare=%q", base, compare)
	}
}

func TestResolveRefs_ThreeDots(t *testing.T) {
	// "main...feature" should be treated as range with "main" and "..feature".
	base, compare := ResolveRefs([]string{"main...feature"}, "", "")
	if base != "main" || compare != ".feature" {
		t.Errorf("got base=%q compare=%q", base, compare)
	}
}

func TestResolveRefs_FlagsOverride(t *testing.T) {
	base, compare := ResolveRefs([]string{"ignored"}, "develop", "staging")
	if base != "develop" || compare != "staging" {
		t.Errorf("expected flags to override, got base=%q compare=%q", base, compare)
	}
}

func TestResolveRefs_HeadTilde(t *testing.T) {
	base, compare := ResolveRefs([]string{"HEAD~3"}, "", "")
	if base != "HEAD~3" {
		t.Errorf("expected base=HEAD~3, got %q", base)
	}
	if compare != "" {
		t.Errorf("expected empty compare, got %q", compare)
	}
}

func TestGitInstalled(t *testing.T) {
	// Git should be installed in CI and dev environments.
	if !GitInstalled() {
		t.Skip("git not installed")
	}
}
