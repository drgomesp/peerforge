package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	abciclient "github.com/tendermint/tendermint/abci/client"
	"github.com/tendermint/tendermint/abci/types"
	tendermintconfig "github.com/tendermint/tendermint/config"
	tendermintlog "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/libs/service"
	tendermintnode "github.com/tendermint/tendermint/node"
	"github.com/urfave/cli/v2"

	"github.com/drgomesp/peerforge/internal/peerforge-hubd/abci"
	_ "github.com/drgomesp/peerforge/internal/peerforge-hubd/abci"
)

//var socketAddr string
var configFile string

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	//flag.StringVar(&socketAddr, "socket-addr", "tcp://0.0.0.0:26658", "Unix domain socket address")
	home := os.Getenv("HOME")
	flag.StringVar(&configFile, "config", fmt.Sprintf("%s/%s", home, "/.tendermint/config/config.toml"),
		"Path to config.toml")
}

func main() {
	app := &cli.App{
		Name: "📡 peerforge-hubd",
		Action: func(context *cli.Context) error {
			app := abci.NewApplication()

			flag.Parse()

			node, err := newTendermint(app, configFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
				os.Exit(2)
			}

			node.Start()
			defer func() {
				node.Stop()
				node.Wait()
			}()

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			<-c
			//
			//flag.Parse()
			//
			//app := abci.NewApplication()
			//logger := tendermint.MustNewDefaultLogger(
			//	tendermint.LogFormatPlain,
			//	tendermint.LogLevelDebug,
			//	false,
			//)
			//
			//server := abciserver.NewSocketServer(socketAddr, app)
			//server.SetLogger(logger)
			//
			//if err := server.Start(); err != nil {
			//	return err
			//}
			//
			//defer server.Stop()
			//
			//c := make(chan os.Signal, 1)
			//signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			//<-c
			//

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Err(err).Send()
	}
}

func newTendermint(app types.Application, configFile string) (service.Service, error) {
	// read config
	config := tendermintconfig.DefaultValidatorConfig()
	config.SetRoot(filepath.Dir(filepath.Dir(configFile)))

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper failed to read config file: %w", err)
	}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("viper failed to unmarshal config: %w", err)
	}
	if err := config.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("config is invalid: %w", err)
	}

	// create logger
	logger, err := tendermintlog.NewDefaultLogger(tendermintlog.LogFormatPlain, config.LogLevel, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// create node
	node, err := tendermintnode.New(
		config,
		logger,
		abciclient.NewLocalCreator(app),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new Tendermint node: %w", err)
	}

	return node, nil
}