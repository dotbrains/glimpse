# __PROJECT_NAME__ — __PROJECT_DESCRIPTION__

![__PROJECT_NAME__](./assets/og-image.svg)

[![CI](https://github.com/dotbrains/__PROJECT_NAME__/actions/workflows/ci.yml/badge.svg)](https://github.com/dotbrains/__PROJECT_NAME__/actions/workflows/ci.yml)
[![Release](https://github.com/dotbrains/__PROJECT_NAME__/actions/workflows/release.yml/badge.svg)](https://github.com/dotbrains/__PROJECT_NAME__/actions/workflows/release.yml)
[![License: PolyForm Shield 1.0.0](https://img.shields.io/badge/License-PolyForm%20Shield%201.0.0-blue.svg)](https://polyformproject.org/licenses/shield/1.0.0/)

![Go](https://img.shields.io/badge/-Go-00ADD8?style=flat-square&logo=go&logoColor=white)
![Cobra](https://img.shields.io/badge/-Cobra-00ADD8?style=flat-square&logo=go&logoColor=white)
![macOS](https://img.shields.io/badge/-macOS-000000?style=flat-square&logo=apple&logoColor=white)
![Linux](https://img.shields.io/badge/-Linux-FCC624?style=flat-square&logo=linux&logoColor=black)

__PROJECT_DESCRIPTION_LONG__

## Problem

<!-- Describe the problem this CLI solves. Why does it exist? What's wrong with existing approaches? -->

## Configuration

`__PROJECT_NAME__` reads its configuration from `~/.config/__PROJECT_NAME__/config.yaml`. If the file does not exist, built-in defaults are used.

### Config file format

```yaml
# Add your config schema here.
# Example:
# default_agent: claude-cli
#
# agents:
#   claude-cli:
#     provider: claude-cli
#     model: opus
```

### `__PROJECT_NAME__ config init`

Scaffolds a config file with the built-in defaults:

```
$ __PROJECT_NAME__ config init
✓ Wrote default config to ~/.config/__PROJECT_NAME__/config.yaml
Edit the file to customize settings.
```

Refuses to overwrite an existing file unless `--force` is passed.

## Commands

### `__PROJECT_NAME__`

<!-- Describe the main command behavior here. -->

Steps:
1. Step one.
2. Step two.
3. Step three.

```
$ __PROJECT_NAME__
→ output example here
```

## Architecture

```
main.go                           Entry point, version injection via ldflags
cmd/
  root.go                         Cobra root command + subcommand registration
  cmd_test.go                     Command-level tests
internal/
  config/
    config.go                     YAML config: Load, Save, defaults
    config_test.go
  exec/
    executor.go                   CommandExecutor interface for testability
    executor_test.go
```

### Adding new subcommands

1. Create `cmd/<name>.go` with a `newXxxCmd()` factory function.
2. Register it in `root.go` via `root.AddCommand(newXxxCmd())`.
3. Add tests in `cmd/cmd_test.go`.

### Adding internal packages

1. Create a new directory under `internal/<package>/`.
2. Include `*_test.go` files alongside source files.
3. Use the `exec.CommandExecutor` interface for any shell-out logic.

## Testing

```sh
# Run all tests with race detection
make test

# Run with coverage report
make cover

# Lint
make lint
```

Tests use `t.TempDir()` and `t.Setenv()` for isolation. The `exec.MockExecutor` pattern is used to test code that shells out to external commands.

## Release

Releases are triggered by pushing a git tag:

```sh
git tag v0.1.0
git push origin v0.1.0
```

This triggers the release workflow which:
1. Runs tests and lint.
2. Builds via GoReleaser for darwin/linux × amd64/arm64.
3. Publishes a GitHub release with binaries.
4. Updates the Homebrew tap at `dotbrains/homebrew-tap`.
