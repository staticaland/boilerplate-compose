package main

import (
    "flag"
    "fmt"
    "log"
    "os"

    "github.com/staticaland/boilerplate-compose/internal/bpcompose"
)

var (
    version = "0.1.0"
)

func main() {
    var configFiles multiStringFlag
    var varFiles multiStringFlag
    var varPairs multiStringFlag
    var printConfig bool
    var showVersion bool

    flag.Var(&configFiles, "f", "Path to a boilerplate-compose file. May be specified multiple times. If not set, defaults to searching boilerplate-compose.yaml/yml in current directory.")
    flag.Var(&varFiles, "var-file", "Load variable values from the YAML file FILE. May be specified multiple times.")
    flag.Var(&varPairs, "var", "Set variable NAME=VALUE. May be specified multiple times.")
    flag.BoolVar(&printConfig, "print-config", false, "Print the fully resolved configuration and exit.")
    flag.BoolVar(&showVersion, "version", false, "Print version and exit.")

    flag.Parse()

    if showVersion {
        fmt.Println(version)
        return
    }

    loader := bpcompose.NewLoader()

    // Determine starting files
    startFiles := make([]string, 0)
    if len(configFiles) > 0 {
        startFiles = append(startFiles, configFiles...)
    } else {
        // auto-detect
        if _, err := os.Stat("boilerplate-compose.yaml"); err == nil {
            startFiles = append(startFiles, "boilerplate-compose.yaml")
        } else if _, err := os.Stat("boilerplate-compose.yml"); err == nil {
            startFiles = append(startFiles, "boilerplate-compose.yml")
        } else {
            log.Fatal("no boilerplate-compose.yaml or boilerplate-compose.yml found; specify with -f")
        }
    }

    resolved, err := loader.LoadAndResolve(startFiles)
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    // Merge var-files and var pairs into top-level vars
    if len(varFiles) > 0 {
        if err := bpcompose.MergeVarFiles(resolved, varFiles); err != nil {
            log.Fatalf("failed to merge var-files: %v", err)
        }
    }
    if len(varPairs) > 0 {
        if err := bpcompose.MergeVarPairs(resolved, varPairs); err != nil {
            log.Fatalf("failed to merge vars: %v", err)
        }
    }

    if printConfig {
        if err := bpcompose.PrintResolved(resolved, os.Stdout); err != nil {
            log.Fatalf("failed to print config: %v", err)
        }
        return
    }

    // For now, we just validate and exit
    if err := bpcompose.Validate(resolved); err != nil {
        log.Fatalf("configuration invalid: %v", err)
    }

    fmt.Println("configuration valid")
}

type multiStringFlag []string

func (m *multiStringFlag) String() string {
    return fmt.Sprintf("%v", []string(*m))
}

func (m *multiStringFlag) Set(value string) error {
    *m = append(*m, value)
    return nil
}