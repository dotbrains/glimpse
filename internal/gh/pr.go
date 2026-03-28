package gh

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/dotbrains/glimpse/internal/comments"
	gexec "github.com/dotbrains/glimpse/internal/exec"
)

var prURLPattern = regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+)/pull/(\d+)`)

// PRInfo holds metadata about a GitHub pull request.
type PRInfo struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	State     string `json:"state"`
	BaseRef   string `json:"baseRefName"`
	HeadRef   string `json:"headRefName"`
	RepoSlug  string `json:"-"`
	Owner     string `json:"-"`
	Repo      string `json:"-"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

// IsPRURL returns true if the string looks like a GitHub PR URL.
func IsPRURL(s string) bool {
	return prURLPattern.MatchString(s)
}

// ParsePRURL extracts owner, repo, and PR number from a GitHub PR URL.
func ParsePRURL(url string) (owner, repo string, number string, err error) {
	m := prURLPattern.FindStringSubmatch(url)
	if m == nil {
		return "", "", "", fmt.Errorf("not a valid GitHub PR URL: %s", url)
	}
	return m[1], m[2], m[3], nil
}

// Client wraps GitHub CLI operations.
type Client struct {
	Executor gexec.CommandExecutor
}

// NewClient creates a new GitHub CLI client.
func NewClient() *Client {
	return &Client{Executor: gexec.NewRealExecutor()}
}

// GHInstalled returns true if gh is on PATH.
func GHInstalled() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

// FetchPRDiff fetches the unified diff for a PR.
func (c *Client) FetchPRDiff(ctx context.Context, owner, repo, number string) (string, error) {
	out, err := c.Executor.Run(ctx, "gh", "pr", "diff", number, "-R", owner+"/"+repo)
	if err != nil {
		return "", fmt.Errorf("fetching PR diff: %w", err)
	}
	return out, nil
}

// FetchPRInfo fetches metadata for a PR.
func (c *Client) FetchPRInfo(ctx context.Context, owner, repo, number string) (*PRInfo, error) {
	fields := "number,title,body,state,baseRefName,headRefName,additions,deletions"
	out, err := c.Executor.Run(ctx, "gh", "pr", "view", number, "-R", owner+"/"+repo, "--json", fields)
	if err != nil {
		return nil, fmt.Errorf("fetching PR info: %w", err)
	}
	var info PRInfo
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &info); err != nil {
		return nil, fmt.Errorf("parsing PR info: %w", err)
	}
	info.Owner = owner
	info.Repo = repo
	info.RepoSlug = owner + "/" + repo
	return &info, nil
}

// PRComment represents a GitHub PR review comment.
type PRComment struct {
	Path string `json:"path"`
	Line int    `json:"line"`
	Body string `json:"body"`
	Side string `json:"side,omitempty"` // "RIGHT" or "LEFT"
}

// PostReviewComments posts comments as a PR review via gh api.
func (c *Client) PostReviewComments(ctx context.Context, owner, repo, number string, cs []comments.Comment) error {
	if len(cs) == 0 {
		return nil
	}

	var ghComments []PRComment
	for _, comment := range cs {
		if comment.Resolved {
			continue
		}
		side := "RIGHT"
		if comment.Side == "old" {
			side = "LEFT"
		}
		tag := ""
		if comment.Severity != "" {
			tag = "[" + comment.Severity + "] "
		}
		ghComments = append(ghComments, PRComment{
			Path: comment.File,
			Line: comment.Line,
			Body: tag + comment.Body,
			Side: side,
		})
	}

	if len(ghComments) == 0 {
		return nil
	}

	payload := map[string]interface{}{
		"event":    "COMMENT",
		"comments": ghComments,
	}
	data, _ := json.Marshal(payload)

	_, err := c.Executor.RunWithStdin(ctx, string(data), "gh", "api",
		"repos/"+owner+"/"+repo+"/pulls/"+number+"/reviews",
		"--method", "POST", "--input", "-")
	if err != nil {
		return fmt.Errorf("posting PR review: %w", err)
	}
	return nil
}

// FetchPRComments fetches existing review comments from a PR.
func (c *Client) FetchPRComments(ctx context.Context, owner, repo, number string) ([]comments.Comment, error) {
	out, err := c.Executor.Run(ctx, "gh", "api",
		"repos/"+owner+"/"+repo+"/pulls/"+number+"/comments",
		"--jq", `.[] | {path: .path, line: (.line // .original_line), body: .body, side: .side}`)
	if err != nil {
		return nil, fmt.Errorf("fetching PR comments: %w", err)
	}

	var result []comments.Comment
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line == "" {
			continue
		}
		var raw struct {
			Path string `json:"path"`
			Line int    `json:"line"`
			Body string `json:"body"`
			Side string `json:"side"`
		}
		if json.Unmarshal([]byte(line), &raw) != nil {
			continue
		}
		side := "new"
		if raw.Side == "LEFT" {
			side = "old"
		}
		result = append(result, comments.Comment{
			File:   raw.Path,
			Line:   raw.Line,
			Side:   side,
			Body:   raw.Body,
			Author: "github",
		})
	}
	return result, nil
}
