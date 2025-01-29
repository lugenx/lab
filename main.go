package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const labVersion = "1.3.0"

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
		i := strings.Index(line, "=")
		if i != -1 {
			key := strings.TrimSpace(line[:i])
			value := strings.TrimSpace(line[i+1:])
			if strings.Contains(value, "#") {
				fmt.Printf(Yellow+"\n  Comments are not supported in config. Please remove '#' and everything after it from: "+Reset+Red+"%s"+Reset+Yellow+"in %s.lab\n\n"+Reset, line, displayPath)
				os.Exit(1)
			}
			config[key] = value
		}
	}

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
		OpenFile(labdir, firstArg, config["editor"], displayPath)
	} else {
		CreateAndOpenFile(labdir, config["prefix"], firstArg, config["editor"], displayPath)
	}
}
