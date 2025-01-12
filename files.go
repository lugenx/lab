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
	dir, _ := os.ReadDir(labdir)

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

	_, err := os.Create(file)
	if err != nil {
		log.Fatalf("failed to create file %v", err)
	}
	cmd := exec.Command(editor, file)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	fileInfo, _ := os.Stat(file)

	if fileInfo.Size() == 0 {
		os.Remove(file)
	}
}

func ListFiles(labdir string) {
	dir, _ := os.ReadDir(labdir)

	sort.Slice(dir, func(i, j int) bool {
		infoI, _ := dir[i].Info()
		infoJ, _ := dir[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})

	for n, file := range dir {
		fileName := file.Name()
		if fileName != ".lab" {
			fmt.Printf("[%2d] %v\n", n+1, fileName)
		}
	}
}

func OpenFile(labdir string, tag string, editor string) {
	dir, _ := os.ReadDir(labdir)

	sort.Slice(dir, func(i, j int) bool {
		infoI, _ := dir[i].Info()
		infoJ, _ := dir[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})

	for n, file := range dir {
		fileName := file.Name()
		file := filepath.Join(labdir, fileName)

		if tag == strconv.Itoa(n+1) {

			cmd := exec.Command(editor, file)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			cmd.Run()
		}
	}
}
