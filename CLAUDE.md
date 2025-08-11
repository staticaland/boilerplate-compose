# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

boilerplate-compose is a CLI tool written in Go that orchestrates multiple boilerplate templates into a single project setup. It reads YAML configuration files defining templates and executes the Gruntwork boilerplate CLI with proper streaming, error handling, and execution reporting.

## Commands

### Building

```bash
go build -o boilerplate-compose
```

### Testing

```bash
# Run all Go tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run specific package tests
go test ./config -v
go test ./processor -v
go test ./executor -v

# Integration tests using Just
just
just test-basic
just test-compose
just test-complex
```

### Running the CLI

```bash
# Basic usage with default config
./boilerplate-compose

# Specify config file
./boilerplate-compose -f my-compose.yaml

# Dry run to preview commands
./boilerplate-compose -dry-run

# Verbose output
./boilerplate-compose -verbose
```

## Architecture

The codebase follows a clean modular architecture:

### Core Packages

1. **main.go** - CLI entry point and flag handling
2. **config/** - YAML configuration parsing and validation
   - `types.go`: Core data structures (ComposeConfig, Template, etc.)
   - `loader.go`: YAML loading with include/extends support
3. **processor/** - Template processing logic
   - `template.go`: Converts Template configs to CLI arguments
   - `orchestrator.go`: Coordinates template execution and provides summary reporting
4. **executor/** - CLI execution with streaming output
   - `cli.go`: Executes boilerplate CLI with real-time output streaming
   - `result.go`: Execution result tracking and summary generation

### Key Data Flow

1. CLI flags parsed in main.go
2. Config file loaded and validated by config package
3. TemplateProcessor converts templates to ProcessingJobs (CLI arguments)
4. Orchestrator coordinates execution via CliExecutor
5. CliExecutor streams boilerplate CLI output in real-time
6. ExecutionSummary provides detailed reporting

### Configuration Structure

Templates support extensive configuration options:

- `template-url`, `output-folder` (required)
- `vars`, `var-file` (string or array)
- Boolean flags: `non-interactive`, `no-hooks`, `no-shell`, `disable-dependency-prompt`
- Actions: `missing-key-action`, `missing-config-action`

Advanced features include file includes and template inheritance via `include` and `extends`.

## Dependencies

- Go 1.24.4+
- `gopkg.in/yaml.v3` for YAML parsing
- External dependency: [Gruntwork boilerplate CLI](https://github.com/gruntwork-io/boilerplate) must be available in PATH

## Testing Strategy

- Unit tests for each package with `*_test.go` files
- Integration tests using Just recipes (Justfile in root, test files in `tests/` directory)
- Test files include `basic-test.yaml` and `complex-test.yaml` for different scenarios
- Dry-run mode for safe testing without actual template execution
