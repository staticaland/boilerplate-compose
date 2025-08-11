package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"boilerplate-compose/config"
	"boilerplate-compose/processor"
)

var (
	configFile = flag.String("f", "", "Path to compose file")
	version    = flag.Bool("version", false, "Show version")
	help       = flag.Bool("help", false, "Show help")
	dryRun     = flag.Bool("dry-run", false, "Show what would be executed without running")
)

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
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

	templateProcessor := processor.NewTemplateProcessor(cfg)
	orchestrator := processor.NewOrchestrator(templateProcessor, *dryRun)

	if err := orchestrator.Process(); err != nil {
		log.Fatalf("Processing failed: %v", err)
	}

	if *dryRun {
		fmt.Println("\nDry run completed. Use without --dry-run to execute.")
	} else {
		fmt.Println("All templates processed successfully.")
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