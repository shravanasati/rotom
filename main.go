package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
	"os"
)

var errDownloadFirst = errors.New("run the `download` command first")

func main() {
	rotom := &cli.Command{
		Name:  "rotom",
		Usage: "Show pokemon image on terminal.",
		Description: "Run it without any arguments to get random pokemon image. You can pass a dex/name to show its image too.",
		Action: func(ctx context.Context, c *cli.Command) error {
			dexOrName := c.Args().Get(0)
			if dexOrName == "" {
				filename, err := getRandomFile(SPRITES_DIR)
				if err != nil {
					return fmt.Errorf("unable to find any pokemon in SPRITES_DIR=%s, %w", SPRITES_DIR, errDownloadFirst)
				}
				return DisplayImage(filename)

			} else {
				filename, err := searchFile(dexOrName)
				if err != nil {
					return err
				}
				return DisplayImage(filename)
			}
		},

		Commands: []*cli.Command{
			{
				Name:  "download",
				Usage: "Download all pokemon sprites.",
				Action: func(ctx context.Context, c *cli.Command) error {
					DownloadAllSprites()
					return nil
				},
			},
		},
	}

	if err := rotom.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err)
	}
}
