package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
   lab <extension>                   Create and open a new file (e.g., lab js)
   lab <number>                      Open file by number
   lab 0                             Open config file
   lab                               List all lab files
   lab -v, --version                 Show version
   lab -h, --help                    Show this help
   lab -d, --delete <number>         Delete file by number
   lab -p, --path                    Show file path
   lab -r, --run <number> <command>  Run any command on the specified file (e.g., python, node, cat)

 Examples:
   lab js                  Create a JavaScript file
   lab 1                   Open most recent file
   lab 0                   Edit config
   lab                     List all lab files
   lab -d 2                Delete file #2
   lab -p 1                Show the path of file #1
   lab -r 1 node           Run the file #1 with Node.js
   lab -r 2 vim            Open in different editor
   lab -r 3 cat            View file contents

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
			fmt.Printf("\n  " + Yellow + "No file specified for deletion. Use 'lab -d <number>'\n\n" + Reset)
			return
		}
		err := os.Remove(fileDir)
		if err != nil {
			fmt.Printf(Red+"Error: Failed to delete '%s': %v\n"+Reset, file.Name(), err)
			return
		}
		fmt.Printf("\n  "+Red+"%v "+Reset+"has been deleted from the lab!\n\n", file.Name())

	case "-p", "--path":
		if fileDir == "" || file == nil {
			fmt.Printf("\n  " + Yellow + "No file specified for retrieving the path. Use 'lab -p <number>'\n\n" + Reset)
			return
		}
		fmt.Println(fileDir)

	case "-r", "--run":
		if (len(os.Args)) < 4 {
			fmt.Println("\n  " + Yellow + "Missing arguments. Use 'lab -r <number> <command>'\n\n" + Reset)
			return
		}
		runner := strings.Join(os.Args[3:], " ")
		if fileDir == "" || file == nil {
			fmt.Println(Red + "No file specified to run. Provide a valid file number." + Reset)
			return
		}
		var cmd *exec.Cmd
		if strings.Contains(runner, "'") || strings.Contains(runner, "\"") {
			parts := parseCommand(runner)
			cmd = exec.Command(parts[0], append(parts[1:], fileDir)...)
		} else {
			parts := strings.Fields(runner)
			cmd = exec.Command(parts[0], append(parts[1:], fileDir)...)
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Printf(Red+"Error: Failed to execute '%s' on file '%s': %v\n"+Reset, runner, fileDir, err)
		}
	}
}
