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

## Quick Start

```sh
# Install
go install github.com/dotbrains/__PROJECT_NAME__@latest

# Show version
__PROJECT_NAME__ --version

# Initialize config
__PROJECT_NAME__ config init
```

## How It Works

<!-- Describe the core workflow of your CLI here. -->

1. Step one
2. Step two
3. Step three

## Installation

### Via `go install`

```sh
go install github.com/dotbrains/__PROJECT_NAME__@latest
```

### Via Homebrew

```sh
brew tap dotbrains/tap
brew install --cask __PROJECT_NAME__
```

### Via GitHub Release

```sh
gh release download --repo dotbrains/__PROJECT_NAME__ --pattern '__PROJECT_NAME___darwin_arm64.tar.gz' --dir /tmp
tar -xzf /tmp/__PROJECT_NAME___darwin_arm64.tar.gz -C /usr/local/bin
```

### From source

```sh
git clone https://github.com/dotbrains/__PROJECT_NAME__.git
cd __PROJECT_NAME__
make install
```

## Configuration

```sh
# Create default config
__PROJECT_NAME__ config init

# Config lives at ~/.config/__PROJECT_NAME__/config.yaml
```

See [SPEC.md](SPEC.md) for the full config format.

## Commands

| Command | Description |
|---|---|
| `__PROJECT_NAME__` | __PROJECT_DESCRIPTION__ |
| `__PROJECT_NAME__ config init` | Create default config file |

## Dependencies

- **[Go](https://go.dev/)** >= 1.24

## License

This project is licensed under the [PolyForm Shield License 1.0.0](https://polyformproject.org/licenses/shield/1.0.0/) — see [LICENSE](LICENSE) for details.
