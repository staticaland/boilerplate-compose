package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"boilerplate-compose/config"
)

var (
	configFile = flag.String("f", "", "Path to compose file")
	version    = flag.Bool("version", false, "Show version")
	help       = flag.Bool("help", false, "Show help")
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

	fmt.Printf("Loaded config from: %s\n", configPath)
	fmt.Printf("Found %d templates:\n", len(cfg.Templates))
	for name, template := range cfg.Templates {
		fmt.Printf("  - %s: %s -> %s\n", name, template.TemplateURL, template.OutputFolder)
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