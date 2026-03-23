# glimpse — GitHub-style Git Diff Viewer CLI

![glimpse](./assets/og-image.svg)

[![CI](https://github.com/dotbrains/glimpse/actions/workflows/ci.yml/badge.svg)](https://github.com/dotbrains/glimpse/actions/workflows/ci.yml)
[![Release](https://github.com/dotbrains/glimpse/actions/workflows/release.yml/badge.svg)](https://github.com/dotbrains/glimpse/actions/workflows/release.yml)
[![License: PolyForm Shield 1.0.0](https://img.shields.io/badge/License-PolyForm%20Shield%201.0.0-blue.svg)](https://polyformproject.org/licenses/shield/1.0.0/)

![Go](https://img.shields.io/badge/-Go-00ADD8?style=flat-square&logo=go&logoColor=white)
![Cobra](https://img.shields.io/badge/-Cobra-00ADD8?style=flat-square&logo=go&logoColor=white)
![macOS](https://img.shields.io/badge/-macOS-000000?style=flat-square&logo=apple&logoColor=white)
![Linux](https://img.shields.io/badge/-Linux-FCC624?style=flat-square&logo=linux&logoColor=black)

Browser-based, GitHub-style diff viewer for git changes. View uncommitted changes, branch comparisons, commit ranges, and GitHub PRs with syntax-highlighted split diffs, inline comments, and AI code review.

## Quick Start

```sh
# Install
go install github.com/dotbrains/glimpse@latest

# View uncommitted changes
glimpse

# Compare branches
glimpse main..feature

# View a GitHub PR
glimpse https://github.com/owner/repo/pull/123

# AI code review
glimpse review

# Resolve comments
glimpse resolve
```

## How It Works

1. Run `glimpse` in any git repo — your browser opens with a GitHub-style diff view.
2. Click any line number gutter to leave inline comments with severity tags (`must-fix`, `suggestion`, `nit`, `question`).
3. Run `glimpse review` to get AI-powered code review — comments appear inline in the viewer.
4. Run `glimpse resolve` to output all open comments for your AI agent to fix.

## See Your Diffs

```sh
# Everyday use
glimpse                                    # review all uncommitted changes
glimpse HEAD~1                             # review your last commit
glimpse HEAD~3                             # review your last 3 commits

# Branch workflows
glimpse main                               # compare current branch against main
glimpse main..feature                      # compare feature branch against main
glimpse main feature                       # same as above, shorthand syntax
glimpse --base main --compare feature      # same as above, explicit flags

# Releases and tags
glimpse v1.0.0 v2.0.0                     # compare two releases

# GitHub PRs
glimpse https://github.com/owner/repo/pull/123
```

## AI Code Review

Run `glimpse review` to have an AI agent review the diff and leave inline comments:

```sh
# Review working tree changes
glimpse review

# Review with focus area
glimpse review --focus security
glimpse review --focus performance

# Review specific refs
glimpse review main..feature
```

Comments appear inline in the viewer with severity badges: `[must-fix]`, `[suggestion]`, `[nit]`, `[question]`.

## Resolve Comments

Run `glimpse resolve` to output all open comments for your AI agent to action:

```sh
# Output all open comments
glimpse resolve

# Resolve a specific comment by ID
glimpse resolve abc123
```

A typical workflow: run `glimpse review` to get AI feedback, check the comments in the browser, then run `glimpse resolve` to have your agent apply the fixes.

## Installation

### Via `go install`

```sh
go install github.com/dotbrains/glimpse@latest
```

### Via Homebrew

```sh
brew tap dotbrains/tap
brew install --cask glimpse
```

### Via GitHub Release

```sh
gh release download --repo dotbrains/glimpse --pattern 'glimpse_darwin_arm64.tar.gz' --dir /tmp
tar -xzf /tmp/glimpse_darwin_arm64.tar.gz -C /usr/local/bin
```

### From source

```sh
git clone https://github.com/dotbrains/glimpse.git
cd glimpse
make install
```

## Commands

| Command | Description |
|---|---|
| `glimpse` | View uncommitted changes (staged + unstaged) |
| `glimpse <ref>` | View changes since a ref (branch, tag, commit) |
| `glimpse <base> <compare>` | Compare two refs |
| `glimpse <base>..<compare>` | Range syntax |
| `glimpse <PR_URL>` | View a GitHub PR diff |
| `glimpse review [refs]` | Run AI code review and post inline comments |
| `glimpse review --focus <area>` | Focus review on security, performance, etc. |
| `glimpse resolve [id]` | Output open comments for your agent to fix |
| `glimpse list` | Show all running instances |
| `glimpse list --json` | Machine-readable instance list |
| `glimpse config init` | Create default config file |

## Options

```
--base <ref>       Base ref to compare from (e.g. main, HEAD~3, v1.0.0)
--compare <ref>    Ref to compare against base (default: working tree)
--port <port>      Custom port (default: auto-assigned from 5391)
--no-open          Don't open browser
--quiet            Minimal terminal output
--new              Stop existing instance and start fresh
```

## Dependencies

- **[Go](https://go.dev/)** >= 1.24
- **[git](https://git-scm.com/)**
- **[gh](https://cli.github.com/)** — required for PR URLs
- **[claude](https://docs.anthropic.com/en/docs/claude-code)** — required for `glimpse review` (default AI agent)

## License

This project is licensed under the [PolyForm Shield License 1.0.0](https://polyformproject.org/licenses/shield/1.0.0/) — see [LICENSE](LICENSE) for details.
