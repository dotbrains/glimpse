# Go CLI Template — AI Agent Instructions

This is a template repository for creating Go CLI tools under the `dotbrains` org. It encodes the conventions from `prr`, `awsw`, `aptscout`, `distill`, and `gh-identity`.

## How to Use This Template

### Step 1: Clone and rename

```sh
git clone https://github.com/dotbrains/go-cli-template.git <project-name>
cd <project-name>
rm -rf .git
git init
```

### Step 2: Find-and-replace placeholders

Replace these placeholders across **all files** (including filenames where applicable):

| Placeholder | Replace with | Example |
|---|---|---|
| `__PROJECT_NAME__` | CLI binary name (lowercase, hyphen-separated) | `prr`, `awsw`, `aptscout` |
| `__PROJECT_DESCRIPTION__` | One-line description for `--help` and badges | `AI-powered PR code review CLI` |
| `__PROJECT_DESCRIPTION_LONG__` | 1-2 sentence description for README and SPEC | `Run AI-powered code reviews on GitHub pull requests...` |
| `__MODULE_PATH__` | Go module path (already uses `github.com/dotbrains/__PROJECT_NAME__`) | — |

Quick command:

```sh
# Replace in all files (macOS sed)
find . -type f -not -path './.git/*' -exec sed -i '' \
  -e 's/__PROJECT_NAME__/myproject/g' \
  -e 's/__PROJECT_DESCRIPTION__/My project description/g' \
  -e 's/__PROJECT_DESCRIPTION_LONG__/A longer description of what this tool does./g' \
  {} +
```

### Step 3: Initialize Go module

```sh
rm go.mod go.sum
go mod init github.com/dotbrains/<project-name>
go mod tidy
```

### Step 4: Verify

```sh
make build
make test
```

## What's Included

### Required components (do NOT remove)

| File/Dir | Purpose |
|---|---|
| `main.go` | Entry point with version injection via ldflags |
| `cmd/root.go` | Cobra root command using `newRootCmd(version)` factory pattern |
| `cmd/cmd_test.go` | Baseline tests for Execute, version, subcommands, config init |
| `internal/config/` | YAML config management (Load, Save, DefaultConfig, ConfigDir, ConfigPath) |
| `internal/exec/` | `CommandExecutor` interface + `RealExecutor` for testable shell-outs |
| `Makefile` | Standard targets: build, test, lint, install, clean, vet, cover |
| `.goreleaser.yaml` | GoReleaser config: darwin+linux, amd64+arm64, Homebrew cask tap |
| `.golangci.yml` | golangci-lint v2 config with standard linters |
| `.github/workflows/ci.yml` | CI: test (ubuntu+macos matrix), lint, build |
| `.github/workflows/release.yml` | Release: test+lint+build → GoReleaser on tag push |
| `.gitignore` | Binary, dist, coverage, IDE, OS, website patterns |
| `LICENSE` | PolyForm Shield 1.0.0 license |
| `README.md` | OG image, badges, Quick Start, How It Works, Installation (4 methods), Commands table |
| `SPEC.md` | Detailed specification: Problem, Configuration, Commands, Architecture, Testing, Release |
| `assets/` | Place `og-image.svg` here |
| `website/` | Next.js + Tailwind marketing site (8 section components, dark theme, responsive) |

### Customizing the website

The website uses generic `accent-primary`, `accent-secondary`, `accent-tertiary` color names in `tailwind.config.js`. Rename these to match your project branding (e.g. `prr-amber`, `distill-teal`).

Files to customize:
- `website/tailwind.config.js` — accent color names and hex values
- `website/src/styles/globals.css` — `text-gradient` utility references accent colors
- `website/app/layout.tsx` — metadata (title, description, OG tags)
- `website/src/components/sections/*.tsx` — all section content (features, stats, use cases, code examples)
- `website/public/favicon.svg` and `website/public/og-image.svg` — replace placeholders
- `website/package.json` — update `name`, `description`, dev port

After adding the website, update `Makefile` test target to exclude it: `go test ... $$(go list ./... | grep -v /website/)`.

### Optional components (add as needed)

These are NOT included but should be added for projects that need them:

#### Documentation directory (`docs/`)

For CLIs with complex configuration or shell integrations:

```
docs/
├── architecture.md
├── commands.md
├── configuration.md
├── contributing.md
└── shell-integration.md   # if the CLI has shell wrapper functions
```

#### AI agent integration (`internal/agent/`)

For CLIs that use AI providers (like `prr` and `distill`):

```
internal/agent/
├── agent.go              # Agent interface
├── registry.go           # Provider registry
├── anthropic/            # Claude API provider
├── openai/               # OpenAI API provider
├── claudecli/            # Claude CLI provider
└── codexcli/             # Codex CLI provider
```

Add `agents` subcommand to `cmd/` to list configured agents.

#### Shell integration (`shell/`)

For CLIs that need shell wrapper functions (like `awsw` and `gh-identity`):

```
shell/
├── awsw.fish
├── awsw.bash
└── awsw.zsh
```

Add `init <shell>` subcommand to output the wrapper function.

#### Embedded web UI (`internal/server/`)

For CLIs that serve a local web UI (like `prr` and `aptscout`):

```
internal/server/
├── server.go
└── static/
    ├── index.html
    ├── app.js
    ├── style.css
    └── favicon.svg
```

Add `serve` subcommand to `cmd/`.

## Conventions to Follow

### Code patterns

- **One subcommand per file** in `cmd/` (e.g. `cmd/review.go`, `cmd/history.go`).
- **Factory functions** for commands: `newXxxCmd() *cobra.Command`.
- **`_test.go` files** alongside every source file.
- **`internal/` for all domain logic** — `cmd/` should be thin wrappers.
- **Use `exec.CommandExecutor` interface** for anything that shells out.
- **Use `config.Load()`** to read config — always fall back to `DefaultConfig()`.
- **Hide completion command**: `CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true}`.

### Testing patterns

- Use `t.TempDir()` for temporary directories.
- Use `t.Setenv("HOME", tmp)` to isolate config paths.
- Test via `newRootCmd()` + `root.SetArgs()` + `root.Execute()` for integration-style tests.
- Use `bytes.Buffer` with `root.SetOut(buf)` to capture output.
- Use `MockExecutor` for testing code that shells out.
- Target 80%+ test coverage.

### Documentation patterns

- **README.md** sections in order: title + badges, description, Quick Start, How It Works, Installation (4 methods), Configuration, Commands table, Dependencies, License.
- **SPEC.md** is the detailed spec — everything a developer needs to understand and extend the tool.
- **Badges**: CI, Release, License, then tech stack (Go, Cobra, any APIs, macOS, Linux).

### Release patterns

- Version injected via ldflags: `-X main.version={{.Version}}`.
- GoReleaser builds for darwin+linux × amd64+arm64.
- Homebrew cask published to `dotbrains/homebrew-tap`.
- macOS quarantine removed via post-install hook.
- Changelog excludes `docs:` and `chore:` commits.

## Checklist for New Projects

- [ ] Replaced all `__PROJECT_NAME__` placeholders
- [ ] Replaced all `__PROJECT_DESCRIPTION__` placeholders
- [ ] Replaced all `__PROJECT_DESCRIPTION_LONG__` placeholders
- [ ] Initialized `go.mod` with correct module path
- [ ] Created `assets/og-image.svg`
- [ ] Added project-specific config fields to `internal/config/config.go`
- [ ] Added project-specific subcommands to `cmd/`
- [ ] Added domain logic packages under `internal/`
- [ ] Tests pass with `make test`
- [ ] Lint passes with `make lint`
- [ ] Updated README.md with real Quick Start examples
- [ ] Updated SPEC.md with real config format and command docs
- [ ] Updated Commands table in README.md
- [ ] Added tech stack badges to README.md (e.g. Anthropic, OpenAI, SQLite)
- [ ] Created GitHub repo and pushed
- [ ] Added `HOMEBREW_TAP_TOKEN` secret to GitHub repo settings
- [ ] Removed this `TEMPLATE.md` file
