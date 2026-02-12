package main

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

func getAllFilesDir(dir string) ([]os.DirEntry, error) {
	emptySlice := []os.DirEntry{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return emptySlice, err
	}

	// Filter out directories, keep only files
	var files []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry)
		}
	}

	if len(files) == 0 {
		return emptySlice, fmt.Errorf("no files found in directory")
	}

	return files, nil
}

func getRandomFile(dir string) (string, error) {
	files, err := getAllFilesDir(dir)
	if err != nil {
		return "", err
	}

	// Pick a random file
	randomIndex := rand.IntN(len(files))
	randomFile := files[randomIndex]

	return filepath.Join(dir, randomFile.Name()), nil
}

func searchPokemon(dexOrName string) (string, error) {
	_, err := strconv.Atoi(dexOrName)
	isDex := err == nil
	var searchPattern string
	if isDex {
		searchPattern = fmt.Sprintf("%s-*.png", dexOrName)
	} else {
		searchPattern = fmt.Sprintf("*-%s.png", dexOrName)
	}

	files, err := searchFiles(searchPattern)
	if err != nil {
		return "", fmt.Errorf("fatal: malformed search pattern: %w", err)
	}
	if len(files) == 0 {
		var otherMatchingNames []string
		if !isDex {
			// perform * search only if name is provided
			wildcardSearch := fmt.Sprintf("*%s*", dexOrName)
			otherMatchingNames, err = searchFiles(wildcardSearch)
			if err != nil {
				panic("other matching names search: " + err.Error())
			}
		}

		var didYouMeanText string
		if len(otherMatchingNames) > 0 {
			didYouMeanText = "\n\ndid you mean: " + strings.Join(slices.Collect(Map(slices.Values(otherMatchingNames), pokemonNameFromFilename)), "/") + "?"
		}

		return "", fmt.Errorf("pokemon `%s` not found, ensure correct dex/name and running `download` command%s", dexOrName, didYouMeanText)
	}

	return files[0], nil
}

func Map[T, S any](s iter.Seq[T], f func(i T) S) iter.Seq[S] {
	return func(yield func(S) bool) {
		for v := range s {
			if !yield(f(v)) {
				break
			}
		}
	}
}

func searchFiles(pattern string) ([]string, error) {
	return filepath.Glob(filepath.Join(SPRITES_DIR, pattern))
}
