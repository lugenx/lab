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

// craete file will take labdir and extension and create a file with correct extension in that directory

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

	_, err = os.Create(file)
	if err != nil {
		log.Fatalf("failed to create file %v", err)
	}
	cmd := exec.Command(editor, file)

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

func ListFiles(labdir string, lifedays string) {
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
	)
	// less than 2, because there there is already .lab file
	if len(organizedFiles) < 2 {
		fmt.Printf("\n\t%sYour lab is empty!%s Create a new file with: %slab <extension>%s (e.g., %slab js%s)\n\n",
			Cyan, Reset, Green, Reset, Yellow, Reset)
		return
	}

	fmt.Printf("\n  To open, use: lab" + Green + " <number>\n" + Reset)
	fmt.Printf("  To create: lab <extension>\n\n")

	fmt.Printf("\t\033[36mLab File(s):\033[0m\n\n")
	// fmt.Printf("\t            "+Grey+"%v\n"+Reset, labdir)
	for i, file := range organizedFiles {

		info, _ := file.Info()
		age := time.Since(info.ModTime())
		daysLeft := int(float64(days) - age.Hours()/24)

		fmt.Printf("\t"+Green+"%5s"+Reset+" [%dd]  %v\n", fmt.Sprintf("[%d] ", i+1), daysLeft, file.Name())
	}
	fmt.Println("")
	// fmt.Printf("\n\033[35m  Tip:\033[0m Frequently modified files might be worth keeping permanently\n\n")
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
		cmd := exec.Command(editor, fullFileName)

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
