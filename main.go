package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"boilerplate-compose/config"
	"boilerplate-compose/processor"
	"boilerplate-compose/executor"
)

var (
	configFile      = flag.String("f", "", "Path to compose file")
	version         = flag.Bool("version", false, "Show version")
	help            = flag.Bool("help", false, "Show help")
	dryRun          = flag.Bool("dry-run", false, "Show what would be executed without running")
	boilerplatePath = flag.String("boilerplate-path", "", "Path to boilerplate CLI (defaults to PATH lookup)")
	verbose         = flag.Bool("verbose", false, "Show detailed output from boilerplate commands")
)

func main() {
	flag.Parse()

	if *help {
		printUsage()
		return
	}

	if *version {
		fmt.Println("boilerplate-compose version 0.1.0")
		return
	}

	configPath := findConfigFile(*configFile)
	if configPath == "" {
		log.Fatal("No compose file found. Use -f to specify a file.")
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	templateProcessor := processor.NewTemplateProcessor(cfg, configPath)
	cliExecutor := executor.NewCliExecutor(*boilerplatePath, *verbose)
	orchestrator := processor.NewOrchestrator(templateProcessor, cliExecutor, *dryRun)

	if err := orchestrator.Process(); err != nil {
		log.Fatalf("Processing failed: %v", err)
	}

	if *dryRun {
		fmt.Println("\nDry run completed. Use without --dry-run to execute.")
	} else {
		fmt.Println("\nAll templates processed successfully.")
	}
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

func printUsage() {
	fmt.Println("boilerplate-compose - Orchestrate template rendering using boilerplate CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  boilerplate-compose [options]")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println("\nExample:")
	fmt.Println("  boilerplate-compose -f my-compose.yaml --verbose")
	fmt.Println("  boilerplate-compose --dry-run")
}