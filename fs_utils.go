package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
)

func getRandomFile(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	// Filter out directories, keep only files
	var files []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry)
		}
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no files found in directory")
	}

	// Pick a random file
	randomIndex := rand.IntN(len(files))
	randomFile := files[randomIndex]

	return filepath.Join(dir, randomFile.Name()), nil
}

func searchFile(dexOrName string) (string, error) {
	_, err := strconv.Atoi(dexOrName)
	isDex := err == nil
	var searchPattern string
	if isDex {
		searchPattern = fmt.Sprintf("%s-*.png", dexOrName)
	} else {
		searchPattern = fmt.Sprintf("*-%s.png", dexOrName)
	}

	files, err := filepath.Glob(filepath.Join(SPRITES_DIR, searchPattern))
	if err != nil {
		return "", fmt.Errorf("fatal: malformed search pattern: %w", err)
	}
	if len(files) == 0 {
		return "", fmt.Errorf("pokemon `%s` not found, ensure correct dex/name and running `download` command", dexOrName)
	}

	return files[0], nil
}
