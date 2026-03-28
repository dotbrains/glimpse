# glimpse

![glimpse](./assets/og-image.svg)

[![CI](https://github.com/dotbrains/glimpse/actions/workflows/ci.yml/badge.svg)](https://github.com/dotbrains/glimpse/actions/workflows/ci.yml)
[![Release](https://github.com/dotbrains/glimpse/actions/workflows/release.yml/badge.svg)](https://github.com/dotbrains/glimpse/actions/workflows/release.yml)
[![License: PolyForm Shield 1.0.0](https://img.shields.io/badge/License-PolyForm%20Shield%201.0.0-blue.svg)](https://polyformproject.org/licenses/shield/1.0.0/)

Glimpse is an agent-agnostic, GitHub-style diff viewer and code review tool.

```bash
go install github.com/dotbrains/glimpse@latest
```

It works with Claude Code, Cursor, Codex, and any AI coding agent.

| What can you do? | Description |
|---|---|
| [See your diffs](#see-your-diffs) | View changes in working area, across commits, branches, tags, etc |
| [AI code review](#ai-code-review) | Let your agent review code and leave comments on the diff |
| [Browse project files](#browse-project-files) | Explore your repo and comment on any file for AI to resolve |
| [GitHub PRs](#github-prs) | Pull down a PR, review it locally, push comments back to GitHub |
| [Multiple projects](#multiple-projects) | Run it in multiple repos at once, each gets its own port |

## See your diffs

Run `glimpse` inside any git repo — your browser opens with a GitHub-style, syntax-highlighted diff.

```bash
# everyday use
glimpse                                    # review all uncommitted changes
glimpse HEAD~1                             # review your last commit
glimpse HEAD~3                             # review your last 3 commits

# branch workflows
glimpse main                               # compare current branch against main
glimpse main..feature                      # compare feature branch against main
glimpse main feature                       # same as above, shorthand syntax
glimpse --base main --compare feature      # same as above, explicit flags

# releases and tags
glimpse v1.0.0 v2.0.0                     # compare two releases
glimpse v1.0.0                             # what changed since v1.0.0

# specific commits
glimpse abc1234                            # changes since a specific commit
glimpse abc1234..def5678                   # changes between two commits
```

The `--base`/`--compare` flags use the same terminology as GitHub PRs — base is what you're comparing against, compare is the branch with changes. You can also use range syntax (`main..feature`) or just pass two positional args (`glimpse main feature`).

You can leave comments on any diff — working tree changes, branch comparisons, commit ranges. Copy them into your agent with a button and ask it to resolve them, or use the skills below to let your agent auto-review and auto-solve them.

## AI code review

Install the skills for your coding agent (Claude Code, Cursor, Codex, etc.):

```bash
# Skills are in .skills/ — copy to your agent's skills directory
cp .skills/*.md ~/.claude/skills/   # Claude Code
```

Then use the slash commands:

### `/glimpse-diff`

Opens the diff viewer in your browser. Accepts the same refs as the CLI, plus natural language:

```
/glimpse-diff                          # working tree changes
/glimpse-diff main                     # current branch against main
/glimpse-diff main..feature            # branch diff
/glimpse-diff HEAD~1                   # last commit
/glimpse-diff last 3 commits           # natural language works too
```

Leave comments on any line — when you're done, run `/glimpse-resolve` to have your agent fix them.

### `/glimpse-review`

Your agent reviews the diff and leaves inline comments in the viewer. Uses severity tags (`[must-fix]`, `[suggestion]`, `[nit]`, `[question]`) so you can triage by importance. Supports refs, focus areas, and natural language:

```
/glimpse-review                             # review working tree changes
/glimpse-review main                        # review what you're merging into main
/glimpse-review main..feature               # review branch diff
/glimpse-review identify security issues    # focus on security issues
/glimpse-review performance in src/lib      # focus on performance in specific dir
/glimpse-review last 3 commits              # natural language works too
```

### `/glimpse-resolve`

Reads all open comments and makes the requested code changes. Works with both your comments and AI review comments:

```
/glimpse-resolve                       # resolve all open comments
/glimpse-resolve abc123                # resolve a specific thread by ID
```

A typical workflow: run `/glimpse-review` to get AI feedback, check the comments in the browser, then run `/glimpse-resolve` to apply the fixes.

## Browse project files

Run `glimpse tree` to open a full file tree browser — no diff required. Browse your repo, read files with syntax highlighting, and leave comments on any file or line.

```bash
glimpse tree
```

The tree view supports the same commenting and resolve workflow as the diff viewer. Leave comments on specific lines, files, or folders, then have your agent resolve them.

### `/glimpse-tree`

Opens the file tree browser:

```
/glimpse-tree
```

### `/glimpse-resolve-tree`

Reads open comments from the tree browser and makes the requested code changes:

```
/glimpse-resolve-tree                  # resolve all open comments
/glimpse-resolve-tree abc123           # resolve a specific thread by ID
```

## GitHub PRs

Pass a GitHub PR URL to view and review pull requests locally:

```bash
glimpse https://github.com/owner/repo/pull/123
```

This fetches the PR diff, opens it against its base branch, and lets you leave comments in the viewer. Requires the `gh` CLI installed and authenticated (`gh auth login`).

You can push your comments (including AI review comments) back to GitHub as PR review comments, and pull existing GitHub comments into the viewer. Both are available from the viewer UI.

The skills work with PR URLs too:

```
/glimpse-diff https://github.com/owner/repo/pull/123
/glimpse-review https://github.com/owner/repo/pull/123
```

## Multiple projects

Glimpse supports running multiple projects simultaneously. Each gets its own port automatically:

```bash
# Terminal 1 — starts on :5391
cd ~/projects/app && glimpse

# Terminal 2 — starts on :5392
cd ~/projects/api && glimpse
```

If you run `glimpse` in a repo that already has a running instance, it opens the existing one instead of starting a new server. Use `--new` to kill the existing instance and start fresh.

```bash
glimpse list               # show all running instances
glimpse list --json        # machine-readable output
```

## Options

```
--base <ref>       Base ref to compare from (e.g. main, HEAD~3, v1.0.0)
--compare <ref>    Ref to compare against base (default: working tree)
--port <port>      Custom port (default: auto-assigned from 5391)
--no-open          Don't open browser
--dark             Dark mode (default: true, use --dark=false for light)
--unified          Unified view (default: split)
--quiet            Minimal terminal output
--new              Stop existing instance and start fresh
```

## License

PolyForm Shield 1.0.0 © dotbrains
