# baselith

A cross-platform command-line application built with Go and Cobra that works on Linux, macOS, and Windows.

## Features

- Cross-platform support (Linux, macOS, Windows)
- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Automated testing and building via GitHub Actions
- Example commands demonstrating CLI functionality

## Installation

### From Source

1. Clone the repository:
```bash
git clone https://github.com/hinha/baselith.git
cd baselith
```

2. Build the application:
```bash
go build -o baselith .
```

### From Release

Download the pre-built binary for your platform from the [Releases](https://github.com/hinha/baselith/releases) page.

## Usage

### Basic Commands

Run the application without arguments to see the welcome message:
```bash
./baselith
```

Check the version:
```bash
./baselith version
```

Greet someone:
```bash
./baselith greet -n "World"
# or
./baselith greet --name "World"
```

Get help:
```bash
./baselith --help
./baselith greet --help
```

## Development

### Running Tests

```bash
go test ./... -v
```

### Building for Different Platforms

Linux:
```bash
GOOS=linux GOARCH=amd64 go build -o baselith-linux-amd64 .
```

macOS:
```bash
GOOS=darwin GOARCH=amd64 go build -o baselith-darwin-amd64 .
```

Windows:
```bash
GOOS=windows GOARCH=amd64 go build -o baselith-windows-amd64.exe .
```

## GitHub Actions

This project includes two GitHub Actions workflows:

1. **PR Workflow** (`.github/workflows/pr.yml`): 
   - Runs on pull requests to the main branch
   - Executes unit tests
   - Builds the application for Linux, macOS, and Windows
   - Uploads build artifacts

2. **Release Workflow** (`.github/workflows/release.yml`):
   - Runs when a new release is created
   - Builds the application for all platforms
   - Runs tests on each platform
   - Executes the binary to verify functionality
   - Uploads release assets

## Branch Protection

To enable branch protection for the main branch:

1. Go to repository Settings → Branches
2. Add a branch protection rule for `main`
3. Enable "Require a pull request before merging"
4. Set "Required number of approvals before merging" to at least 1
5. Enable "Require status checks to pass before merging"
6. Select the status checks from the PR workflow

## License

Apache License 2.0 - see [LICENSE](LICENSE) file for details.