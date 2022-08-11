package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/drgomesp/peerforge/internal/peerforge-hubd/abci"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	abciserver "github.com/tendermint/tendermint/abci/server"
	tendermint "github.com/tendermint/tendermint/libs/log"
	"github.com/urfave/cli/v2"
)

var socketAddr string

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	flag.StringVar(&socketAddr, "socket-addr", "tcp://0.0.0.0:26658", "Unix domain socket address")
}

func main() {
	app := &cli.App{
		Name: "📡 peerforge-hubd",
		Action: func(context *cli.Context) error {
			flag.Parse()

			app := abci.NewApplication()
			logger := tendermint.MustNewDefaultLogger(
				tendermint.LogFormatPlain,
				tendermint.LogLevelDebug,
				false,
			)

			server := abciserver.NewSocketServer(socketAddr, app)
			server.SetLogger(logger)

			if err := server.Start(); err != nil {
				return err
			}

			defer server.Stop()

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			<-c

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Err(err).Send()
	}
}
