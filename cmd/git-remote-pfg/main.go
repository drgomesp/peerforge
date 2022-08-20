package main

import (
	"os"
	"strings"

	gitremotepfg "github.com/drgomesp/peerforge/internal/git-remote-pfg"
	"github.com/drgomesp/peerforge/pkg/gitremote"
	shell "github.com/ipfs/go-ipfs-api"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

const IpfsURL = "localhost:45005"
const EmptyRepo = "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// TODO: remove this
	if os.Getenv(shell.EnvDir) == "" {
		_ = os.Setenv(shell.EnvDir, IpfsURL)
	}
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal().Msg("gitremote-remote-pfg expects 2 arguments (origin name and url)")
	}

	remoteName := os.Args[2]
	if strings.HasPrefix(remoteName, "pfg://") {
		remoteName = remoteName[len("pfg://"):]
	}

	if remoteName == "" {
		remoteName = EmptyRepo
	}

	if os.Getenv("GIT_DIR") == "" {
		log.Warn().Msg("missing repository path ($GIT_DIR)... using current directory")
		cwd, err := os.Getwd()
		if err != nil {
			log.Err(err).Send()
		}
		os.Setenv("GIT_DIR", cwd)
	}

	abci, err := rpchttp.New("http://localhost:26657")
	if err != nil {
		log.Err(err).Send()
	}

	handler, err := gitremotepfg.NewPfg(abci, os.Getenv(shell.EnvDir), remoteName)
	if err != nil {
		log.Err(err).Send()
	}

	proto, err := gitremote.NewProtocol("prefix", handler)
	if err != nil {
		log.Err(err).Send()
	}

	if err := proto.Run(os.Stdin, os.Stdout); err != nil {
		log.Err(err).Send()
	}
}
