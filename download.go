package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"

	"github.com/mtslzr/pokeapi-go"
	"github.com/schollz/progressbar/v3"
)

var errSkipDownload = errors.New("skipped download")

const NUM_WORKERS = 20
const MAX_RETRIES = 5

var SPRITES_DIR string
var downloadedSprites []int

type downloadConfig struct {
	maxDex        int
	forceDownload bool
}

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

	files, err := getAllFilesDir(SPRITES_DIR)
	if err != nil {
		return
	}
	downloadedSprites = slices.Collect(Map(slices.Values(files), func(o os.DirEntry) int {
		dex, _ := pokemonFromFilename(o.Name())
		return dex
	}))
}

func DownloadAllSprites(dc downloadConfig) {
	bar := progressbar.Default(int64(dc.maxDex), "sprites downloaded")
	successCount := 0
	newDownloadCount := 0

	jobsChan := make(chan int)
	resultsChan := make(chan error)
	for range NUM_WORKERS {
		go downloadPokemonSpriteWorker(jobsChan, resultsChan, dc.forceDownload)
	}

	// send all jobs
	go func() {
		for i := range dc.maxDex {
			dex := i + 1
			jobsChan <- dex
		}
		close(jobsChan)
	}()

	// collect results
	for range dc.maxDex {
		e := <-resultsChan
		if e == errSkipDownload {
			successCount += 1
		} else if e != nil {
			fmt.Println("error:", e)
		} else {
			newDownloadCount += 1
			successCount += 1
		}
		bar.Add(1)
	}

	fmt.Printf("Downloaded %v sprites (%v new).\n", successCount, newDownloadCount)
}

func downloadPokemonSpriteWorker(jobs chan int, results chan error, forceDownload bool) {
	for dex := range jobs {
		results <- downloadPokemonSprite(dex, forceDownload)
	}
}

func downloadPokemonSprite(dex int, forceDownload bool) error {
	if slices.Contains(downloadedSprites, dex) && !forceDownload {
		// dont redownload
		return errSkipDownload
	}

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
	var lastErr error
	for range MAX_RETRIES {
		resp, err := http.Get(url)
		if err != nil {
			lastErr = err
			continue
		}

		content, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}

		return os.WriteFile(filepath.Join(SPRITES_DIR, filename), content, 0644)
	}

	return lastErr
}
