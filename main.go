package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var errDownloadFirst = errors.New("run the `download` command first")

func main() {
	rotom := &cli.Command{
		Name:        "rotom",
		Usage:       "Show pokemon image on terminal.",
		Description: "Run it without any arguments to get random pokemon image. You can pass a dex/name to show its image too.",
		Action: func(ctx context.Context, c *cli.Command) error {
			dexOrName := NormalizePokemonName(strings.Join(c.Args().Slice(), " "))
			if dexOrName == "" {
				filename, err := getRandomFile(SPRITES_DIR)
				if err != nil {
					return fmt.Errorf("unable to find any pokemon in SPRITES_DIR=%s, %w", SPRITES_DIR, errDownloadFirst)
				}
				return DisplayImage(filename)

			} else {
				filename, err := searchPokemon(dexOrName)
				if err != nil {
					return err
				}
				return DisplayImage(filename)
			}
		},

		Commands: []*cli.Command{
			{
				Name:        "download",
				Usage:       "Download all pokemon sprites.",
				Description: fmt.Sprintf("All sprites are downloaded at `%s`.", SPRITES_DIR),
				Action: func(ctx context.Context, c *cli.Command) error {
					DownloadAllSprites()
					return nil
				},
			},
			{
				Name:  "version",
				Usage: "Print version information.",
				Action: func(ctx context.Context, c *cli.Command) error {
					fmt.Printf("rotom %s (commit: %s, built at: %s)\n", version, commit, date)
					return nil
				},
			},
		},
	}

	if err := rotom.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err)
	}
}
