package main

import (
	"flag"
	"fmt"
	"os"
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

	configPath := *configFile
	if configPath == "" {
		// Look for default files
		if _, err := os.Stat("boilerplate-compose.yaml"); err == nil {
			configPath = "boilerplate-compose.yaml"
		} else if _, err := os.Stat("boilerplate-compose.yml"); err == nil {
			configPath = "boilerplate-compose.yml"
		} else {
			fmt.Println("No compose file found. Use -f to specify a file.")
			os.Exit(1)
		}
	}

	fmt.Printf("Hello World! Using config file: %s\n", configPath)
}