package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const labVersion = "1.1.1"

func checkRequiredConfigs(cfg map[string]string, keys []string) {
	for _, key := range keys {
		if value, ok := cfg[key]; !ok || value == "" {
			fmt.Printf("No %s set, please configure using .lab file inside your /lab directory\n", key)
		}
	}
}

func main() {
	labdir, configFile, displayPath := Setup()
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("failed to read config file %v", err)
	}

	content := string(data)
	contentLines := strings.Split(content, "\n")

	config := make(map[string]string)

	for _, line := range contentLines {
		// TODO: WIP
		// parts := strings.Split(line, "=")
		// fmt.Println("---->", parts)
		// if len(parts) == 2 {
		// 	key := strings.TrimSpace(parts[0])
		// 	value := strings.TrimSpace(parts[1])
		// 	config[key] = value
		// }
		i := strings.Index(line, "=")
		if i != -1 {
			key := strings.TrimSpace(line[:i])
			value := strings.TrimSpace(line[i+1:])
			config[key] = value
		}
	}
	// ---------------------------
	requiredKeys := []string{"editor", "lifedays"}

	checkRequiredConfigs(config, requiredKeys)
	if err := DeleteExpiredFiles(labdir, config["lifedays"]); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-") {
		organizedFiles := organizeFiles(labdir)
		handleFlags(labVersion, organizedFiles, labdir)
		return
	}

	if len(os.Args) == 1 {
		ListFiles(labdir, config["lifedays"], displayPath)
		return
	}

	firstArg := os.Args[1]
	if _, err := strconv.ParseInt(firstArg, 10, 64); err == nil {
		OpenFile(labdir, firstArg, config["editor"])
	} else {
		CreateAndOpenFile(labdir, config["prefix"], firstArg, config["editor"])
	}
}
