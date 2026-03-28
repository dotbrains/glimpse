package git

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	gexec "github.com/dotbrains/glimpse/internal/exec"
)

// Client wraps git operations.
type Client struct {
	Executor gexec.CommandExecutor
	RepoDir  string // working directory for git commands
}

// NewClient creates a Client that operates in the given directory.
func NewClient(dir string) *Client {
	return &Client{
		Executor: gexec.NewRealExecutor(),
		RepoDir:  dir,
	}
}

// run executes a git command in the repo directory.
func (c *Client) run(ctx context.Context, args ...string) (string, error) {
	fullArgs := append([]string{"-C", c.RepoDir}, args...)
	return c.Executor.Run(ctx, "git", fullArgs...)
}

// IsRepo returns true if the directory is inside a git repository.
func (c *Client) IsRepo(ctx context.Context) bool {
	_, err := c.run(ctx, "rev-parse", "--is-inside-work-tree")
	return err == nil
}

// TopLevel returns the root of the git repository.
func (c *Client) TopLevel(ctx context.Context) (string, error) {
	out, err := c.run(ctx, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// CurrentBranch returns the name of the current branch, or "HEAD" if detached.
func (c *Client) CurrentBranch(ctx context.Context) string {
	out, err := c.run(ctx, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "HEAD"
	}
	return strings.TrimSpace(out)
}

// DefaultBranch guesses the default branch (main, master, or first remote HEAD).
func (c *Client) DefaultBranch(ctx context.Context) string {
	// Try symbolic-ref of origin/HEAD.
	out, err := c.run(ctx, "symbolic-ref", "refs/remotes/origin/HEAD")
	if err == nil {
		ref := strings.TrimSpace(out)
		return strings.TrimPrefix(ref, "refs/remotes/origin/")
	}
	// Fallback: check if main or master exist.
	for _, b := range []string{"main", "master"} {
		if _, err := c.run(ctx, "rev-parse", "--verify", b); err == nil {
			return b
		}
	}
	return "main"
}

// RevParse resolves a ref to a commit SHA.
func (c *Client) RevParse(ctx context.Context, ref string) (string, error) {
	out, err := c.run(ctx, "rev-parse", "--verify", ref)
	if err != nil {
		return "", fmt.Errorf("unknown ref %q: %w", ref, err)
	}
	return strings.TrimSpace(out), nil
}

// Diff runs git diff between two refs and returns the unified diff output.
// If compare is empty, diffs against the working tree.
func (c *Client) Diff(ctx context.Context, base, compare string) (string, error) {
	args := []string{"diff", "--no-color"}
	if base != "" && compare != "" {
		args = append(args, base+"..."+compare)
	} else if base != "" {
		args = append(args, base)
	}
	// If both empty, this diffs uncommitted changes (working tree vs index).
	return c.run(ctx, args...)
}

// DiffStaged returns staged (cached) changes.
func (c *Client) DiffStaged(ctx context.Context) (string, error) {
	return c.run(ctx, "diff", "--cached", "--no-color")
}

// HasChanges returns true if there are uncommitted changes (staged or unstaged).
func (c *Client) HasChanges(ctx context.Context) bool {
	out, _ := c.run(ctx, "status", "--porcelain")
	return strings.TrimSpace(out) != ""
}

// ShortSHA returns the short form of a commit SHA.
func (c *Client) ShortSHA(ctx context.Context, ref string) string {
	out, err := c.run(ctx, "rev-parse", "--short", ref)
	if err != nil {
		return ref
	}
	return strings.TrimSpace(out)
}

// RepoName returns the directory name of the repo root.
func (c *Client) RepoName(ctx context.Context) string {
	top, err := c.TopLevel(ctx)
	if err != nil {
		return "unknown"
	}
	return filepath.Base(top)
}

// ListFiles returns all tracked files in the repo at the given ref.
func (c *Client) ListFiles(ctx context.Context, ref string) ([]string, error) {
	if ref == "" {
		ref = "HEAD"
	}
	out, err := c.run(ctx, "ls-tree", "-r", "--name-only", ref)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

// ShowFile returns the contents of a file at the given ref.
func (c *Client) ShowFile(ctx context.Context, ref, path string) (string, error) {
	if ref == "" {
		ref = "HEAD"
	}
	return c.run(ctx, "show", ref+":"+path)
}

// GitInstalled returns true if git is on PATH.
func GitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// ResolveRefs parses user input into base and compare refs.
// Supports: "main..feature", "main feature", single ref, empty (working tree).
func ResolveRefs(args []string, flagBase, flagCompare string) (base, compare string) {
	// Explicit flags take priority.
	if flagBase != "" {
		base = flagBase
		compare = flagCompare
		return
	}

	if len(args) == 0 {
		return "", ""
	}

	// Single arg with range syntax: "main..feature"
	if len(args) == 1 && strings.Contains(args[0], "..") {
		parts := strings.SplitN(args[0], "..", 2)
		return parts[0], parts[1]
	}

	// Two positional args: "main feature"
	if len(args) >= 2 {
		return args[0], args[1]
	}

	// Single ref: diff against it (e.g. "HEAD~3", "main", "v1.0.0")
	return args[0], ""
}
