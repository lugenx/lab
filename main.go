package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func checkRequiredConfigs(cfg map[string]string, keys []string) {
	for _, key := range keys {
		if value, ok := cfg[key]; !ok || value == "" {
			fmt.Printf("No %s set, please configure using .lab file inside your /lab directory\n", key)
		}
	}
}

func main() {
	labdir, configFile := Setup()

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("failed to read config file %v", err)
	}

	content := string(data)
	contentLines := strings.Split(content, "\n")

	config := make(map[string]string)

	for _, line := range contentLines {
		parts := strings.Split(line, "=")

		if len(parts) == 2 {
			config[parts[0]] = parts[1]
		}
	}

	requiredKeys := []string{"editor", "lifedays", "path", "prefix"}

	checkRequiredConfigs(config, requiredKeys)
	// TODO: Arguments will be checked, adjusted, if its numbers should act differently, handle edge cases
	firstArg := os.Args[1]
	if len(firstArg) < 1 {
		fmt.Println("no arguments")
	}
	CreateAndOpenFile(labdir, config["prefix"], firstArg, config["editor"])
	fmt.Printf("First ar is %v", firstArg)
}
