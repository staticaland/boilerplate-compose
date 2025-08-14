package main

import (
	"flag"
	"fmt"
	"os"

	"boilerplate-compose/config"
	"boilerplate-compose/processor"
	"boilerplate-compose/executor"
)

var (
	// Build-time variables set by goreleaser
	version = "dev"
	
	// CLI flags
	configFile      = flag.String("f", "", "Path to compose file")
	showVersion     = flag.Bool("version", false, "Show version")
	help            = flag.Bool("help", false, "Show help")
	dryRun          = flag.Bool("dry-run", false, "Show what would be executed without running")
	boilerplatePath = flag.String("boilerplate-path", "", "Path to boilerplate CLI (defaults to PATH lookup)")
	verbose         = flag.Bool("verbose", false, "Show detailed output from boilerplate commands")
	envFile         = flag.String("env-file", "", "Path to .env file (defaults to .env in current directory)")
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()

	if *help {
		printUsage()
		return nil
	}

	if *showVersion {
		fmt.Printf("boilerplate-compose version %s\n", version)
		return nil
	}

	configPath := findConfigFile(*configFile)
	if configPath == "" {
		return fmt.Errorf("no compose file found. Use -f to specify a file")
	}

	// Set up environment manager
	envManager := config.NewEnvironmentManager()
	
	// Load system environment first
	envManager.LoadSystemEnvironment()
	
	// Load from .env file if specified or if default .env exists
	envFilePath := findEnvFile(*envFile)
	if envFilePath != "" {
		if err := envManager.LoadEnvironmentFromFile(envFilePath); err != nil {
			return fmt.Errorf("failed to load environment file: %w", err)
		}
	}

	cfg, err := config.LoadConfigWithEnvironment(configPath, envManager)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	templateProcessor := processor.NewTemplateProcessor(cfg, configPath)
	cliExecutor := executor.NewCliExecutor(*boilerplatePath, *verbose)
	orchestrator := processor.NewOrchestrator(templateProcessor, cliExecutor, *dryRun)

	if err := orchestrator.Process(); err != nil {
		return fmt.Errorf("processing failed: %w", err)
	}

	if *dryRun {
		fmt.Println("\nDry run completed. Use without -dry-run to execute.")
	} else {
		fmt.Println("\nAll templates processed successfully.")
	}
	
	return nil
}

func findConfigFile(specified string) string {
	if specified != "" {
		return specified
	}

	candidates := []string{"boilerplate-compose.yaml", "boilerplate-compose.yml"}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return ""
}

func findEnvFile(specified string) string {
	if specified != "" {
		return specified
	}

	// Check for default .env file
	if _, err := os.Stat(".env"); err == nil {
		return ".env"
	}

	return ""
}

func printUsage() {
	fmt.Println("boilerplate-compose - Orchestrate template rendering using boilerplate CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  boilerplate-compose [options]")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println("\nExample:")
	fmt.Println("  boilerplate-compose -f my-compose.yaml -verbose")
	fmt.Println("  boilerplate-compose -dry-run")
	fmt.Println("  boilerplate-compose -env-file production.env")
}