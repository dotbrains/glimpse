package git

import (
	"context"
	"os"
	"os/exec"
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
	if !GitInstalled() {
		t.Skip("git not installed")
	}
}

// initTestRepo creates a temp git repo with one commit and returns the path.
func initTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, args := range [][]string{
		{"init"},
		{"config", "user.email", "test@test.com"},
		{"config", "user.name", "Test"},
	} {
		cmd := append([]string{"-C", dir}, args...)
		out, err := exec.Command("git", cmd...).CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %s %v", args, out, err)
		}
	}
	// Create initial commit.
	os.WriteFile(dir+"/main.go", []byte("package main\n"), 0o644)
	exec.Command("git", "-C", dir, "add", ".").Run()
	exec.Command("git", "-C", dir, "commit", "-m", "init").Run()
	return dir
}

func TestClient_IsRepo(t *testing.T) {
	if !GitInstalled() {
		t.Skip("git not installed")
	}
	repo := initTestRepo(t)
	c := NewClient(repo)
	ctx := context.Background()

	if !c.IsRepo(ctx) {
		t.Error("expected IsRepo = true")
	}

	c2 := NewClient(t.TempDir())
	if c2.IsRepo(ctx) {
		t.Error("expected IsRepo = false for non-repo")
	}
}

func TestClient_TopLevel(t *testing.T) {
	if !GitInstalled() {
		t.Skip("git not installed")
	}
	repo := initTestRepo(t)
	c := NewClient(repo)

	top, err := c.TopLevel(context.Background())
	if err != nil {
		t.Fatalf("TopLevel error: %v", err)
	}
	if top == "" {
		t.Error("expected non-empty top level")
	}
}

func TestClient_CurrentBranch(t *testing.T) {
	if !GitInstalled() {
		t.Skip("git not installed")
	}
	repo := initTestRepo(t)
	c := NewClient(repo)

	branch := c.CurrentBranch(context.Background())
	// Should be master or main depending on git config.
	if branch == "" {
		t.Error("expected non-empty branch")
	}
}

func TestClient_Diff_WorkingTree(t *testing.T) {
	if !GitInstalled() {
		t.Skip("git not installed")
	}
	repo := initTestRepo(t)
	c := NewClient(repo)
	ctx := context.Background()

	// Modify a file.
	os.WriteFile(repo+"/main.go", []byte("package main\n\nfunc main() {}\n"), 0o644)

	diff, err := c.Diff(ctx, "", "")
	if err != nil {
		t.Fatalf("Diff error: %v", err)
	}
	if diff == "" {
		t.Error("expected non-empty diff for modified file")
	}
}

func TestClient_DiffStaged(t *testing.T) {
	if !GitInstalled() {
		t.Skip("git not installed")
	}
	repo := initTestRepo(t)
	c := NewClient(repo)
	ctx := context.Background()

	// Stage a change.
	os.WriteFile(repo+"/new.go", []byte("package main\n"), 0o644)
	exec.Command("git", "-C", repo, "add", "new.go").Run()

	diff, err := c.DiffStaged(ctx)
	if err != nil {
		t.Fatalf("DiffStaged error: %v", err)
	}
	if diff == "" {
		t.Error("expected non-empty staged diff")
	}
}

func TestClient_HasChanges(t *testing.T) {
	if !GitInstalled() {
		t.Skip("git not installed")
	}
	repo := initTestRepo(t)
	c := NewClient(repo)
	ctx := context.Background()

	// Clean repo — no changes.
	if c.HasChanges(ctx) {
		t.Error("expected no changes in clean repo")
	}

	// Modify file.
	os.WriteFile(repo+"/main.go", []byte("package main\n\n// changed\n"), 0o644)
	if !c.HasChanges(ctx) {
		t.Error("expected changes after modifying file")
	}
}

func TestClient_RevParse(t *testing.T) {
	if !GitInstalled() {
		t.Skip("git not installed")
	}
	repo := initTestRepo(t)
	c := NewClient(repo)

	sha, err := c.RevParse(context.Background(), "HEAD")
	if err != nil {
		t.Fatalf("RevParse error: %v", err)
	}
	if len(sha) < 7 {
		t.Errorf("expected SHA, got %q", sha)
	}

	_, err = c.RevParse(context.Background(), "nonexistent-ref-abc")
	if err == nil {
		t.Error("expected error for nonexistent ref")
	}
}

func TestClient_ShortSHA(t *testing.T) {
	if !GitInstalled() {
		t.Skip("git not installed")
	}
	repo := initTestRepo(t)
	c := NewClient(repo)

	short := c.ShortSHA(context.Background(), "HEAD")
	if len(short) < 4 || len(short) > 12 {
		t.Errorf("expected short SHA, got %q", short)
	}
}

func TestClient_RepoName(t *testing.T) {
	if !GitInstalled() {
		t.Skip("git not installed")
	}
	repo := initTestRepo(t)
	c := NewClient(repo)

	name := c.RepoName(context.Background())
	if name == "" || name == "unknown" {
		t.Errorf("expected non-empty repo name, got %q", name)
	}
}

func TestClient_ListFiles(t *testing.T) {
	if !GitInstalled() {
		t.Skip("git not installed")
	}
	repo := initTestRepo(t)
	c := NewClient(repo)

	files, err := c.ListFiles(context.Background(), "HEAD")
	if err != nil {
		t.Fatalf("ListFiles error: %v", err)
	}
	if len(files) == 0 {
		t.Error("expected at least one file")
	}
	found := false
	for _, f := range files {
		if f == "main.go" {
			found = true
		}
	}
	if !found {
		t.Error("expected main.go in file list")
	}
}

func TestClient_ShowFile(t *testing.T) {
	if !GitInstalled() {
		t.Skip("git not installed")
	}
	repo := initTestRepo(t)
	c := NewClient(repo)

	content, err := c.ShowFile(context.Background(), "HEAD", "main.go")
	if err != nil {
		t.Fatalf("ShowFile error: %v", err)
	}
	if content == "" {
		t.Error("expected non-empty file content")
	}

	_, err = c.ShowFile(context.Background(), "HEAD", "nonexistent.go")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestClient_DefaultBranch(t *testing.T) {
	if !GitInstalled() {
		t.Skip("git not installed")
	}
	repo := initTestRepo(t)
	c := NewClient(repo)

	branch := c.DefaultBranch(context.Background())
	if branch == "" {
		t.Error("expected non-empty default branch")
	}
}

func TestClient_Diff_WithRefs(t *testing.T) {
	if !GitInstalled() {
		t.Skip("git not installed")
	}
	repo := initTestRepo(t)
	c := NewClient(repo)
	ctx := context.Background()

	// Add a second commit.
	os.WriteFile(repo+"/second.go", []byte("package main\n"), 0o644)
	exec.Command("git", "-C", repo, "add", ".").Run()
	exec.Command("git", "-C", repo, "commit", "-m", "second").Run()

	diff, err := c.Diff(ctx, "HEAD~1", "")
	if err != nil {
		t.Fatalf("Diff error: %v", err)
	}
	if diff == "" {
		t.Error("expected non-empty diff between commits")
	}
}
