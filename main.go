package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
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

	if err := DeleteExpiredFiles(labdir, config["lifedays"]); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	requiredKeys := []string{"editor", "lifedays", "path", "prefix"}

	checkRequiredConfigs(config, requiredKeys)

	if len(os.Args) == 1 {
		ListFiles(labdir, config["lifedays"])
		return
	}

	firstArg := os.Args[1]

	if _, err := strconv.ParseInt(firstArg, 10, 64); err == nil {
		OpenFile(labdir, firstArg, config["editor"])
	} else {
		CreateAndOpenFile(labdir, config["prefix"], firstArg, config["editor"])
	}
}
