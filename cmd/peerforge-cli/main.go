package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	cli2 "github.com/drgomesp/peerforge/internal/peerforge-cli"
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
					return cli2.Init(ctx.Args().Get(0))
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Err(err).Send()
	}
}