package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/urfave/cli/v2"

	"github.com/peerforge/peerforge/internal/peerforge-cli/repository"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	app := &cli.App{
		Name: "📡 peerforge-cli",
		Commands: []*cli.Command{
			{
				Name:      "init",
				Aliases:   []string{"i"},
				Usage:     `Initializes a project at a given directory`,
				ArgsUsage: "[dir]",
				Action: func(ctx *cli.Context) error {
					abciClient, err := rpchttp.New("http://localhost:26657")
					if err != nil {
						return err
					}
					initializer := repository.NewInitializer(abciClient)
					return initializer.Init(ctx.Args().Get(0))
				},
			},
			{
				Name: "t",
				Action: func(context *cli.Context) error {

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Err(err).Send()
	}
}
