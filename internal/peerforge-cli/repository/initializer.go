package repository

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	ssi "github.com/nuts-foundation/go-did"
	"github.com/nuts-foundation/go-did/did"
	peerforgeevent "github.com/peerforge/peerforge/internal/git-remote-pfg"
	peerforge "github.com/peerforge/peerforge/pkg"
	"github.com/peerforge/peerforge/pkg/gitremote"
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
	return &Initializer{
		abci: abci,
	}
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

	dir, err = gitremote.GetLocalDir()
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

		_, pubKey, err := crypto.GenerateEd25519Key(rand.Reader)
		keyPair, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

		id, err := peer.IDFromPublicKey(pubKey)
		if err != nil {
			return err
		}

		didID, err := did.ParseDID(fmt.Sprintf("did:pfg:%s", id))
		if err != nil {
			return err
		}

		doc := &did.Document{
			ID: *didID,
		}

		keyID, _ := did.ParseDIDURL(fmt.Sprintf("did:pfg:%s#key-1", id))
		verificationMethod, err := did.NewVerificationMethod(*keyID, ssi.JsonWebKey2020, did.DID{}, keyPair.Public())

		// This adds the method to the VerificationMethod list and stores a reference to the assertion list
		doc.AddAssertionMethod(verificationMethod)

		didJson, _ := json.MarshalIndent(doc, "", "  ")

		type EventsTx struct {
			Events []*peerforge.Event `json:"events"`
		}

		data, err := json.Marshal(EventsTx{Events: []*peerforge.Event{
			peerforge.NewEvent(
				peerforgeevent.RepositoryInitialized,
				uuid.New().String(),
				1,
				fmt.Sprintf("did:pfg:%s", id),
			),
		}})

		res, err := i.abci.BroadcastTxCommit(context.Background(), data)
		if err != nil {
			return err
		}
		if res.CheckTx.IsErr() || res.DeliverTx.IsErr() {
			log.Debug().Msgf(err.Error())
			log.Err(err).Send()
			return err
		}

		log.Debug().Msgf(string(didJson))
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

	log.Info().Msgf("Repository initialized.")
	log.Info().Msgf("Push your changes to the PeerForge remote: gitremote push peerforge {branch}")

	return nil
}
