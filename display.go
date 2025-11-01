package main

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/blacktop/go-termimg"
)

func DisplayImage(imagePath string) error {
	termimg.PrintFile(imagePath)
	fmt.Println(pokemonNameFromFilename(imagePath))
	return nil
}

func pokemonNameFromFilename(filename string) string {
	baseName := filepath.Base(filename)
	pokemonNameRegex := regexp.MustCompile(`^\d+-(.+).png$`)
	matches := pokemonNameRegex.FindStringSubmatch(baseName)
	return (matches[1])
}
