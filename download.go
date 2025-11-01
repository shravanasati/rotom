package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/mtslzr/pokeapi-go"
	"github.com/schollz/progressbar/v3"
)

const MAX_DEX = 1025
const NUM_WORKERS = 20

var SPRITES_DIR string

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic("unable to get user home dir: " + err.Error())
	}
	switch runtime.GOOS {
	case "linux", "darwin":
		// Use XDG cache directory on Linux/macOS
		cacheDir := os.Getenv("XDG_CACHE_HOME")
		if cacheDir == "" {
			cacheDir = filepath.Join(homedir, ".cache")
		}
		SPRITES_DIR = filepath.Join(cacheDir, "rotom", "sprites")
	case "windows":
		// Use local app data on Windows
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			localAppData = filepath.Join(homedir, "AppData", "Local")
		}
		SPRITES_DIR = filepath.Join(localAppData, "rotom", "sprites")
	}

	// Create the sprites directory if it doesn't exist
	err = os.MkdirAll(SPRITES_DIR, 0755)
	if err != nil {
		panic("unable to create sprites directory: " + err.Error())
	}
}

func DownloadAllSprites() {
	bar := progressbar.Default(MAX_DEX, "sprites downloaded")

	jobsChan := make(chan int)
	resultsChan := make(chan error)
	for range NUM_WORKERS {
		go downloadPokemonSpriteWorker(jobsChan, resultsChan)
	}

	// send all jobs
	go func() {
		for i := range MAX_DEX {
			dex := i + 1
			jobsChan <- dex
		}
		close(jobsChan)
	}()

	// collect results
	for range MAX_DEX {
		e := <-resultsChan
		if e != nil {
			fmt.Println("error:", e)
		}
		bar.Add(1)
	}
}

func downloadPokemonSpriteWorker(jobs chan int, results chan error) {
	for dex := range jobs {
		results <- downloadPokemonSprite(dex)
	}
}

func downloadPokemonSprite(dex int) error {
	pok, err := pokeapi.Pokemon(strconv.Itoa(dex))
	if err != nil {
		return err
	}

	spriteURL := pok.Sprites.FrontDefault

	filename := fmt.Sprintf("%d-%s.png", dex, pok.Name)
	err = downloadImage(filename, spriteURL)
	return err
}

func downloadImage(filename, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(SPRITES_DIR, filename), content, 0644)
}
