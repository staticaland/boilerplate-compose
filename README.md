# boilerplate-compose

A CLI tool for composing multiple boilerplate templates into a single project setup. Define multiple templates in a YAML configuration file and generate them all at once.

## Features

- **YAML Configuration**: Define multiple templates in a single configuration file
- **Template Variables**: Support for template variables and variable files
- **Real Execution**: Execute boilerplate CLI with streaming output and error handling
- **Dry Run Mode**: Preview commands without executing them
- **Execution Reporting**: Detailed execution summaries with timing and success/failure counts
- **Flexible Configuration**: Support for includes, extends, and various template options
- **Validation**: Built-in configuration validation with clear error messages
- **Auto-discovery**: Automatically finds `boilerplate-compose.yaml` or `boilerplate-compose.yml` files

## Prerequisites

Go 1.19 or later.

[Boilerplate CLI](https://github.com/gruntwork-io/boilerplate) installed and available in PATH. With [Mise](https://mise.jdx.dev/dev-tools/backends/aqua.html):

```
mise use -g aqua:gruntwork-io/boilerplate
```

## Installation

```bash
go build -o boilerplate-compose
```

## Usage

### Basic Usage

```bash
# Use default config file (boilerplate-compose.yaml or boilerplate-compose.yml)
./boilerplate-compose

# Specify a config file
./boilerplate-compose -f my-compose.yaml

# Dry run - preview commands without executing
./boilerplate-compose -dry-run

# Verbose output - show detailed boilerplate CLI output
./boilerplate-compose -verbose

# Custom boilerplate CLI path
./boilerplate-compose -boilerplate-path /usr/local/bin/boilerplate

# Show version
./boilerplate-compose -version

# Show help
./boilerplate-compose -help
```

### Command Line Options

- `-f`: Path to compose configuration file
- `-dry-run`: Show what commands would be executed without running them
- `-verbose`: Show detailed output from boilerplate CLI commands
- `-boilerplate-path`: Path to boilerplate CLI executable (defaults to PATH lookup)
- `-version`: Show version information
- `-help`: Show help message

### Configuration File

Create a `boilerplate-compose.yaml` file:

```yaml
templates:
  frontend:
    template-url: "https://github.com/example/react-template"
    output-folder: "./frontend"
    vars:
      project_name: "my-app"
      author: "john-doe"
    non-interactive: true

  backend:
    template-url: "https://github.com/example/go-api-template"
    output-folder: "./backend"
    var-file: "backend-vars.yaml"
    missing-key-action: "error"
```

### Template Configuration Options

Each template supports the following options:

- `template-url` (required): URL to the template repository
- `output-folder` (required): Where to generate the template
- `vars`: Key-value pairs for template variables
- `var-file`: Path to YAML file with variables (can be string or array)
- `non-interactive`: Skip interactive prompts
- `missing-key-action`: Action when template variables are missing ("error", "skip", etc.)
- `missing-config-action`: Action when config is missing
- `no-hooks`: Disable template hooks
- `no-shell`: Disable shell execution
- `disable-dependency-prompt`: Skip dependency installation prompts

### Advanced Configuration

```yaml
templates:
  app:
    template-url: "https://github.com/example/app-template"
    output-folder: "./app"
    var-file:
      - "common-vars.yaml"
      - "app-specific-vars.yaml"

# Include other compose files
include:
  - path: "shared-templates.yaml"

# Extend from another compose file
extends:
  file: "base-compose.yaml"
  template: "base-template"
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run specific package tests
go test ./config -v
```

### Project Structure

```
.
├── main.go                    # CLI entry point
├── config/
│   ├── types.go              # Configuration data structures
│   ├── loader.go             # YAML parsing and validation
│   ├── types_test.go         # Type tests
│   └── loader_test.go        # Loader tests
├── processor/
│   ├── template.go           # Template processing logic
│   ├── orchestrator.go       # Template orchestration
│   ├── template_test.go      # Template tests
│   └── orchestrator_test.go  # Orchestrator tests
├── executor/
│   ├── cli.go                # CLI execution with streaming
│   ├── result.go             # Execution result tracking
│   ├── cli_test.go           # CLI executor tests
│   └── result_test.go        # Result tests
├── example-compose.yaml       # Example configuration
└── boilerplate-compose.yaml  # Default config file
```

### Execution Output

When running templates, boilerplate-compose provides detailed execution reporting:

```
2025/08/11 16:55:46 Processing 3 templates
2025/08/11 16:55:46 Processing template: frontend
2025/08/11 16:55:46 Processing template: backend
2025/08/11 16:55:46 Processing template: docs

=== Execution Summary ===
Total templates: 3
Successful: 3
Failed: 0
Total duration: 47.75µs

Template execution times:
  ✓ frontend: 28.708µs
  ✓ backend: 6.958µs
  ✓ docs: 8.333µs

All templates processed successfully.
```

### Dry Run Output

Use `-dry-run` to preview what commands will be executed:

```
=== Template: frontend ===
Command that would be executed:
boilerplate --template-url https://github.com/example/react-template --output-folder ./frontend --var project_name=my-react-app --var author=john-doe --var version=1.0.0 --non-interactive --missing-key-action zero

Template details:
  URL: https://github.com/example/react-template
  Output: ./frontend
  Variables:
    project_name = my-react-app
    author = john-doe
    version = 1.0.0

=== Execution Summary ===
Total templates: 3
Successful: 3
Failed: 0
Total duration: 47.75µs

Template execution times:
  ✓ frontend: 28.708µs
  ✓ backend: 6.958µs
  ✓ docs: 8.333µs

Dry run completed. Use without -dry-run to execute.
```

### Error Handling

- **CLI Validation**: Checks if boilerplate CLI is available before execution
- **Stop on Failure**: Execution stops on first template failure (unless in dry-run mode)
- **Clear Error Messages**: Detailed error reporting with template context
- **Execution Summary**: Shows which templates succeeded or failed
