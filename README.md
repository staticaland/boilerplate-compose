# boilerplate-compose

A CLI tool for composing multiple boilerplate templates into a single project setup. Define multiple templates in a YAML configuration file and generate them all at once.

## Features

- **YAML Configuration**: Define multiple templates in a single configuration file
- **Template Variables**: Support for template variables and variable files
- **Flexible Configuration**: Support for includes, extends, and various template options
- **Validation**: Built-in configuration validation with clear error messages
- **Auto-discovery**: Automatically finds `boilerplate-compose.yaml` or `boilerplate-compose.yml` files

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

# Show version
./boilerplate-compose --version

# Show help
./boilerplate-compose --help
```

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
├── example-compose.yaml       # Example configuration
└── boilerplate-compose.yaml  # Default config file
```
