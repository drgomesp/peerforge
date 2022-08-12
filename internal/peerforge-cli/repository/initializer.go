package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	gitremotego "github.com/drgomesp/git-remote-go"
	peerforge "github.com/drgomesp/peerforge/pkg"
	"github.com/drgomesp/peerforge/pkg/event"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tendermint/tendermint/rpc/client"
)

var (
	ErrRepositoryAlreadyInitialized = errors.New("repository already initialized")
)

const (
	RemoteName     = "peerforge"
	ConfigFileName = ".peerforge.yaml"
	DefaultConfig  = `{"foo": "bar"}`
)

type Initializer struct {
	abci client.ABCIClient
}

func NewInitializer(abci client.ABCIClient) *Initializer {
	return &Initializer{abci: abci}
}

// Init initializes an empty Peerforge repository at a given
// directory or existing repository directory.
func (i *Initializer) Init(dir string) (err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if dir == "" {
		dir = cwd
	} else {
		dir = path.Join(cwd, dir)
	}

	_ = os.Setenv("GIT_DIR", dir)
	log.Debug().Str("GIT_DIR", dir).Send()

	dir, err = gitremotego.GetLocalDir()
	if err != nil {
		return err
	}

	log.Info().Msgf("Initializing Peerforge 📡  %s", dir)

	r, err := git.PlainOpen(dir)
	if err != nil {
		if err != git.ErrRepositoryNotExists {
			return err
		}
	}

	if r == nil {
		if r, err = git.PlainInit(dir, false); err != nil {
			return err
		}
	}

	headRef, err := r.Head()
	if err != nil {
		if err != plumbing.ErrReferenceNotFound {
			return err
		}
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(dir, ConfigFileName)); errors.Is(err, os.ErrNotExist) {
		filename := filepath.Join(dir, ConfigFileName)
		err = os.WriteFile(filename, []byte(DefaultConfig), 0644)
		if err != nil {
			return err
		}

		_, err = w.Add(fmt.Sprintf("%s", ConfigFileName))
		if err != nil {
			return err
		}

		log.Info().Msgf("Generated configuration file '%s'", ConfigFileName)
		log.Info().Msgf("Configured 'origin' remote (pfg://)")

		commit, err := w.Commit("initialized Peerforge 📡 repository", &git.CommitOptions{
			Author: &object.Signature{
				Name: "hubd",
				When: time.Now(),
			},
		})
		if err != nil {
			return err
		}

		obj, err := r.CommitObject(commit)
		if err != nil {
			_ = w.Reset(&git.ResetOptions{Commit: headRef.Hash()})
			return err
		}

		log.Debug().Msgf(obj.String())
	} else {
		return ErrRepositoryAlreadyInitialized
	}

	remoteName := RemoteName
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{"pfg://"},
	})
	if err != nil {
		if headRef != nil {
			_ = w.Reset(&git.ResetOptions{Commit: headRef.Hash()})
		}
		return err
	}

	type EventsTx struct {
		Events []*peerforge.Event `json:"events"`
	}

	data, err := json.Marshal(EventsTx{Events: []*peerforge.Event{
		peerforge.NewEvent(
			event.RepositoryInitialized,
			uuid.New().String(),
			1,
			"peerforge.hubd",
		),
	}})
	if err != nil {
		if headRef != nil {
			_ = w.Reset(&git.ResetOptions{Commit: headRef.Hash()})
		}
		return err
	}

	res, err := i.abci.BroadcastTxCommit(context.Background(), data)
	if err != nil {
		if headRef != nil {
			_ = w.Reset(&git.ResetOptions{Commit: headRef.Hash()})
		}
		return err
	}

	if res.CheckTx.IsErr() || res.DeliverTx.IsErr() {
		if headRef != nil {
			_ = w.Reset(&git.ResetOptions{Commit: headRef.Hash()})
		}
		return err
	}

	return nil
}
