package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/blacktop/go-termimg"
)

func DisplayImage(imagePath string) error {
	termimg.PrintFile(imagePath)
	dex, name := pokemonFromFilename(imagePath)
	fmt.Printf("%s (%s)\n", name, formatGeneration(generationFromDex(dex)))
	return nil
}

func pokemonFromFilename(filename string) (int, string) {
	baseName := filepath.Base(filename)
	pokemonNameRegex := regexp.MustCompile(`^(\d+)-(.+).png$`)
	matches := pokemonNameRegex.FindStringSubmatch(baseName)
	dex, err := strconv.Atoi(matches[1])
	if err != nil {
		panic(err)
	}
	return dex, matches[2]
}

func pokemonNameFromFilename(file string) string {
	_, n := pokemonFromFilename(file)
	return n
}
