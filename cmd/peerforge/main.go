package main

import (
	"os"

	"github.com/drgomesp/peerforge/internal/peerforged"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	app := &cli.App{
		Name: "📡 pfgcli",
		Commands: []*cli.Command{
			{
				Name:      "init",
				Aliases:   []string{"i"},
				Usage:     `Initializes a project at a given directory`,
				ArgsUsage: "[dir]",
				Action: func(ctx *cli.Context) error {
					return peerforged.Init(ctx.Args().Get(0))
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Err(err).Send()
	}
}
