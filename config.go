package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func isExists(dirPath string, confPath string) (dir bool, conf bool) {
	_, dirErr := os.Stat(dirPath)
	if dirErr != nil && !os.IsNotExist(dirErr) {
		log.Fatalf("failed to access path %v", dirErr)
	}
	_, confErr := os.Stat(confPath)
	if confErr != nil && !os.IsNotExist(confErr) {
		log.Fatalf("failed to access config %v", confErr)
	}

	return !os.IsNotExist(dirErr), !os.IsNotExist(confErr)
}

//go:embed .lab
var configTemplate string

func Setup() (string, string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("failed to get user home directory", err)
	}

	configDirectory := homeDir

	customLabPath := os.Getenv("LABPATH")

	fmt.Println("custom path here", customLabPath)
	if customLabPath != "" {

		if strings.HasPrefix(customLabPath, "~") {
			customLabPath = strings.Replace(customLabPath, "~", homeDir, 1)
		}

		configDirectory = customLabPath
	}

	labDir := filepath.Join(configDirectory, "lab")
	confFile := filepath.Join(labDir, ".lab")

	hasDir, hasConf := isExists(labDir, confFile)

	if !hasDir {
		err := os.MkdirAll(labDir, 0o755)
		if err != nil {
			log.Fatalf("failed to create directory %v", err)
		}
	}

	if !hasConf {
		newConfigFile, err := os.Create(confFile)
		if err != nil {
			log.Fatalf("failed to create config file %v", err)
		}
		defer newConfigFile.Close()

		_, err = newConfigFile.Write([]byte(configTemplate))
		if err != nil {
			log.Fatalf("failed to write to config file", err)
		}

	}

	return labDir, confFile
}
