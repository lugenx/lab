package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// generate a letter combination for a given count (a, b, ..., z, aa, ab, ..., zzzzz, ...)
func generateLetterCombination(count int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"
	base := len(alphabet)
	combination := ""

	// Convert count to base-26 using letters
	for count >= 0 {
		combination = string(alphabet[count%base]) + combination
		count = count/base - 1
	}

	return combination
}

func parseCommand(cmd string) []string {
	var args []string
	var current string
	var quoteChar rune
	inQuotes := false

	for _, r := range cmd {
		switch r {
		case '\'', '"':
			if inQuotes && r == quoteChar {
				inQuotes = false
			} else if !inQuotes {
				inQuotes = true
				quoteChar = r
			} else {
				current += string(r)
			}
		case ' ':
			if !inQuotes {
				if current != "" {
					args = append(args, current)
					current = ""
				}
			} else {
				current += string(r)
			}
		default:
			current += string(r)
		}
	}
	if current != "" {
		args = append(args, current)
	}
	return args
}

func CreateAndOpenFile(labdir string, prefix string, extension string, editor string) {
	// 2006-01-02 15:04:05
	dir, err := os.ReadDir(labdir)
	if err != nil {
		log.Fatalf("failed to read directory %v", err)
	}

	today := time.Now().Format("060102")
	todaysFilesCount := 0

	for _, file := range dir {
		fileName := file.Name()
		if strings.Contains(fileName, today) {
			todaysFilesCount++
		}
	}

	letterCombination := generateLetterCombination(todaysFilesCount)

	validExtension := extension
	if len(extension) > 0 && extension[:1] == "." {
		validExtension = extension[1:]
	}
	validPrefix := prefix
	if len(prefix) > 0 {
		validPrefix = prefix + "-"
	}

	filename := validPrefix + today + letterCombination + "." + validExtension

	file := filepath.Join(labdir, filename)

	createdFile, err := os.Create(file)
	if err != nil {
		log.Fatalf("failed to create file %v", err)
	}
	createdFile.Close()

	var cmd *exec.Cmd
	// using 'if' to save some performance for users who has single editor command
	if len(strings.Fields(editor)) > 1 {
		parts := parseCommand(editor)
		cmd = exec.Command(parts[0], append(parts[1:], file)...)
	} else {
		cmd = exec.Command(editor, file)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			fmt.Printf("\n  Editor %s not found. \n  \033[33mSet your preferred editor in ~/lab/.lab\033[0m (examples below):\n\n\teditor=code    # for VS Code\n\teditor=nvim    # for Neovim\n\teditor=vim     # for Vim\n\n", editor)
			os.Remove(file)
			return
		}
	}

	fileInfo, _ := os.Stat(file)

	if fileInfo.Size() == 0 {
		os.Remove(file)
	}
}

func organizeFiles(labdir string) []os.DirEntry {
	dir, err := os.ReadDir(labdir)
	if err != nil {
		log.Fatalf("failed to read directory %v", err)
	}
	var noDotLabDir []os.DirEntry
	for _, file := range dir {
		if file.Name() != ".lab" {
			noDotLabDir = append(noDotLabDir, file)
		}
	}
	sort.Slice(noDotLabDir, func(i, j int) bool {
		infoI, _ := noDotLabDir[i].Info()
		infoJ, _ := noDotLabDir[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})
	return noDotLabDir
}

func ListFiles(labdir string, lifedays string, displayPath string) {
	days, err := strconv.ParseFloat(strings.TrimSpace(lifedays), 64)
	if err != nil {
		log.Fatal("invalid lifedays value")
	}

	organizedFiles := organizeFiles(labdir)

	const (
		Reset   = "\033[0m"
		Green   = "\033[32m"
		Cyan    = "\033[36m"
		Magenta = "\033[35m"
		Yellow  = "\033[33m"
		Grey    = "\033[90m"
		Red     = "\033[31m"
	)
	// less than 2, because there there is already .lab file
	if len(organizedFiles) == 0 {
		fmt.Printf("\n\t%sYour lab is empty!%s Create a new file with: %slab <extension>%s (e.g., %slab js%s)\n\n",
			Cyan, Reset, Green, Reset, Yellow, Reset)
		return
	}

	fmt.Printf("\n  To open, use: lab" + Green + " <number>\n" + Reset)
	fmt.Printf("  To create: lab <extension>\n\n")

	fmt.Println("\t" + Cyan + " Lab Files:" + Reset + Grey + "  " + displayPath + "\n" + Reset)
	for i, file := range organizedFiles {

		info, _ := file.Info()
		age := time.Since(info.ModTime())
		var timeLeft string

		daysLeft := int(float64(days) - age.Hours()/24)
		hoursLeft := int(float64(days)*24 - age.Hours())
		minutesLeft := int(float64(days)*24*60 - age.Minutes())

		if minutesLeft < 60 {
			minutesLeftStr := strconv.Itoa(minutesLeft)
			timeLeft = fmt.Sprintf(Red+" %sm"+Reset, minutesLeftStr)
		} else if hoursLeft < 25 {
			hoursLeftStr := strconv.Itoa(hoursLeft)
			timeLeft = fmt.Sprintf(Yellow+" %sh"+Reset, hoursLeftStr)
		} else {
			daysLeftStr := strconv.Itoa(daysLeft)
			timeLeft = fmt.Sprintf(" %sd", daysLeftStr)
		}
		padding := 18
		fmt.Printf("\t"+Green+"%5s "+Reset+"%-*s"+Grey+"%s\n", fmt.Sprintf("[%d] ", i+1), padding, file.Name(), timeLeft)
	}
	fmt.Println("")
}

func OpenFile(labdir string, tag string, editor string) {
	organizedFiles := organizeFiles(labdir)

	if n, err := strconv.Atoi(tag); err == nil && n <= len(organizedFiles) {

		var fileName string

		if n, _ := strconv.Atoi(tag); n == 0 {
			fileName = ".lab"
		} else {
			fileName = organizedFiles[n-1].Name()
		}
		fullFileName := filepath.Join(labdir, fileName)

		var cmd *exec.Cmd
		// using 'if' to save some performance for users who has single editor command
		if len(strings.Fields(editor)) > 1 {

			parts := parseCommand(editor)
			cmd = exec.Command(parts[0], append(parts[1:], fullFileName)...)
		} else {
			cmd = exec.Command(editor, fullFileName)
		}

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			if strings.Contains(err.Error(), "executable file not found") {
				fmt.Printf("\n  Editor %s not found. \n  \033[33mSet your preferred editor in ~/lab/.lab\033[0m (examples below):\n\n\teditor=code    # for VS Code\n\teditor=nvim    # for Neovim\n\teditor=vim     # for Vim\n\n", editor)
				return
			}
		}
	} else {
		fmt.Printf("Invalid file index: %s (index out of range)\n", tag)
	}
}

func DeleteExpiredFiles(labdir string, lifedays string) error {
	organizedFiles := organizeFiles(labdir)
	days, err := strconv.ParseFloat(strings.TrimSpace(lifedays), 64)
	if err != nil {
		return fmt.Errorf("failed to convert string to number %v", err)
	}
	duration := time.Duration(days * 24 * float64(time.Hour))

	for _, file := range organizedFiles {
		info, _ := file.Info()

		if time.Since(info.ModTime()) > duration {
			os.Remove(filepath.Join(labdir, file.Name()))
		}
	}
	return nil
}
