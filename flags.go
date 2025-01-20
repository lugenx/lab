package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const (
	Reset   = "\033[0m"
	Green   = "\033[32m"
	Cyan    = "\033[36m"
	Magenta = "\033[35m"
	Yellow  = "\033[33m"
	Grey    = "\033[90m"
	Red     = "\033[31m"
)

const helpText = `
 Usage:
   lab <extension>            Create and open a new file (e.g., lab js)
   lab <number>               Open file by number
   lab 0                      Open config file
   lab -v, --version          Show version
   lab -h, --help             Show this help
   lab -d, --delete <number>  Delete file by number

 Examples:
   lab js                  Create a JavaScript file
   lab 1                   Open most recent file
   lab 0                   Edit config
   lab -d 2                Delete file #2

 Configuration (~/.lab):
   editor=nvim             Your preferred editor
   lifedays=7              Days to keep files
   prefix=lab              Prefix for filenames

 Files are stored in ~/lab/ (or custom location via LABPATH)
 Files expire after configured days (default 7)
`

func handleFlags(labVersion string, organizedFiles []os.DirEntry, labdir string) {
	if len(os.Args) < 2 {
		fmt.Println("Error: No flag provided. Use -v or -d with an index.")
		return
	}

	flag := os.Args[1]
	var file os.DirEntry
	var fileDir string

	if len(os.Args) > 2 {
		index, err := strconv.Atoi(os.Args[2])
		if err != nil || index < 1 || index > len(organizedFiles) {
			fmt.Println("Invalid or missing file index.")
			return
		}
		file = organizedFiles[index-1]
		fileDir = filepath.Join(labdir, file.Name())
	}

	switch flag {
	case "-v", "--version":
		fmt.Printf("lab version %v\n", labVersion)
		os.Exit(0)
	case "-h", "--help":
		fmt.Println(helpText)
	case "-d", "--delete":
		if fileDir == "" || file == nil {
			fmt.Printf("\n  " + Yellow + "No file specified for deletion.\n\n" + Reset)
			return
		}
		os.Remove(fileDir)
		fmt.Printf("\n  "+Red+"%v "+Reset+"has been deleted from the lab!\n\n", file.Name())
	}
}
