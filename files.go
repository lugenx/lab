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

	cmd.Run()

	fileInfo, _ := os.Stat(file)

	if fileInfo.Size() == 0 {
		os.Remove(file)
	}
}

func ListFiles(labdir string, lifedays string) {
	days, err := strconv.ParseFloat(strings.TrimSpace(lifedays), 64)
	if err != nil {
		log.Fatal("invalid lifedays value")
	}

	dir, err := os.ReadDir(labdir)
	if err != nil {
		log.Fatalf("failed to read directory %v", err)
	}
	sort.Slice(dir, func(i, j int) bool {
		infoI, _ := dir[i].Info()
		infoJ, _ := dir[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})

	for n, file := range dir {
		fileName := file.Name()

		if fileName == ".lab" {
			continue
		}
		info, _ := file.Info()
		age := time.Since(info.ModTime())
		daysLeft := int(float64(days) - age.Hours()/24)
		fmt.Printf("[%2d] [%dd] %v\n", n+1, daysLeft, file.Name())
	}
}

func OpenFile(labdir string, tag string, editor string) {
	dir, err := os.ReadDir(labdir)
	if err != nil {
		log.Fatalf("failed to read directory %v", err)
	}

	sort.Slice(dir, func(i, j int) bool {
		infoI, _ := dir[i].Info()
		infoJ, _ := dir[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})

	// for n, file := range dir {
	// 	fileName := file.Name()
	// 	file := filepath.Join(labdir, fileName)
	//
	// 	if tag == strconv.Itoa(n+1) {
	//
	// 		cmd := exec.Command(editor, file)
	//
	// 		cmd.Stdin = os.Stdin
	// 		cmd.Stdout = os.Stdout
	// 		cmd.Stderr = os.Stderr
	//
	// 		cmd.Run()
	// 	}
	// }

	if n, err := strconv.Atoi(tag); err == nil && n > 0 && n <= len(dir) {
		file := dir[n-1]

		fileName := file.Name()

		fullFileName := filepath.Join(labdir, fileName)
		cmd := exec.Command(editor, fullFileName)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.Run()
	}
}

func DeleteExpiredFiles(labdir string, lifedays string) error {
	files, err := os.ReadDir(labdir)
	if err != nil {
		log.Fatalf("failed to read directory %v", err)
	}

	days, err := strconv.ParseFloat(strings.TrimSpace(lifedays), 64)
	if err != nil {
		return fmt.Errorf("failed to convert string to number %v", err)
	}
	duration := time.Duration(days * 24 * float64(time.Hour))

	for _, file := range files {
		info, _ := file.Info()

		if file.Name() == ".lab" {
			continue
		}

		if time.Since(info.ModTime()) > duration {
			os.Remove(filepath.Join(labdir, file.Name()))
		}
	}
	return nil
}
