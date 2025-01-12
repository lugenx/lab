package main

import (
	_ "embed"
	"os"
)

func hasConfigFile(path string) bool {
	dir, _ := os.Open(path)
	entries, _ := dir.ReadDir(0)

	for _, entry := range entries {
		if entry.Name() == ".lab" {
			return true
		}
		return false
	}

	return false
}

//go:embed .lab
var configTemplate string

func ensureConfigFile() {
	configDirectory, _ := os.UserHomeDir()

	if !hasConfigFile(configDirectory) {
		newConfigFile, _ := os.Create(".lab")
		newConfigFile.Write([]byte(configTemplate))
	}
}
