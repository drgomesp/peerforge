// Package gitremotepfg implements a ProtocolHandler that
// stores git objects as IPLD-structured data in IPFS.
package gitremotepfg

import (
	"bytes"
	"context"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/storage"
	"github.com/rs/zerolog/log"

	gitremote "github.com/drgomesp/git-remote-ipldprime"
	"github.com/drgomesp/git-remote-ipldprime/core"
	ipldgitprime "github.com/drgomesp/go-ipld-gitprime"
	"github.com/drgomesp/go-ipld-gitprime/store"
)

const (
	LargeObjectDir    = "objects"
	LObjTrackerPrefix = "//lobj"
	HEAD              = "HEAD"
)

const (
	RefPathHead = iota
	RefPathRef
)

const (
	RepositoryInitialized = "repository.Initialized"
)

type refPath struct {
	path  string
	rType int

	hash string
}

var _ gitremote.ProtocolHandler = &Pfg{}

type Pfg struct {
	linkSys *ipld.LinkSystem
	repo    *git.Repository
	store   ipldgitprime.Store

	tracker     *core.Tracker
	pushed      bool
	localDir    string
	remoteName  string
	currentHash string
}

func NewPfg(tracker *core.Tracker, remoteName string) (*Pfg, error) {
	cwd, _ := os.Getwd()

	localDir, err := gitremote.GetLocalDir()
	if localDir == "" {
		localDir = cwd
	}

	repo, err := git.PlainOpen(localDir)
	if err == git.ErrWorktreeNotProvided {
		repoRoot, _ := path.Split(localDir)

		repo, err = git.PlainOpen(repoRoot)
		if err != nil {
			return nil, err
		}
	}

	ls := cidlink.DefaultLinkSystem()
	st, err := store.NewObjectStore(&ls)
	ls.SetWriteStorage(st)
	ls.SetReadStorage(st)
	if err != nil {
		return nil, err
	}

	return &Pfg{tracker: tracker, linkSys: &ls, store: st, repo: repo, remoteName: remoteName}, nil
}

func (p *Pfg) Initialize(tracker *core.Tracker, repo *git.Repository) error {
	p.repo = repo
	p.currentHash = p.remoteName

	localDir, err := gitremote.GetLocalDir()
	if err != nil {
		return err
	}

	repo, err = git.PlainOpen(localDir)
	if err != nil {
		return err
	}

	p.localDir = localDir
	p.repo = repo

	return nil
}

func (p *Pfg) Finish() error {
	return nil
}
func (p *Pfg) ProvideBlock(identifier string, tracker *core.Tracker) ([]byte, error) {
	return nil, nil
}

func (p *Pfg) Capabilities() string {
	return gitremote.DefaultCapabilities
}

func (p *Pfg) List(forPush bool) ([]string, error) {
	return []string{}, nil
}

func (p *Pfg) Push(ctx context.Context, local string, remote string) (string, error) {
	p.pushed = true

	localRef, err := p.repo.Reference(plumbing.ReferenceName(local), true)
	if err != nil {
		return "", err
	}

	headHash := localRef.Hash().String()

	push := core.NewPush(p.localDir, p.tracker, p.linkSys, p.repo, p.store)
	push.NewNode = p.bigNodePatcher(p.tracker)

	err = push.PushHash(headHash)
	if err != nil {
		return "", err
	}

	c, err := core.CidFromHex(headHash)
	if err != nil {
		return "", err
	}

	k := path.Join(p.currentHash, remote)
	err = p.store.Put(ctx, k, bytes.NewBufferString(headHash).Bytes())

	head, err := p.getRef(ctx, HEAD)
	if err != nil {
		return "", err
	}

	if head == nil {
		if err = p.store.Put(
			ctx,
			path.Join(p.currentHash, HEAD),
			bytes.NewBufferString("refs/heads/main").Bytes(),
		); err != nil {
			return "", err
		}

		k = path.Join(c.String(), HEAD)
		if err = p.store.Put(ctx, k, c.Bytes()); err != nil {
			return "", err
		}

		p.currentHash = k
	}

	return local, nil
}

func (p *Pfg) bigNodePatcher(tracker *core.Tracker) func(context.Context, string, []byte) error {
	return func(ctx context.Context, hash string, data []byte) error {
		log.Debug().Msgf("size: %vb", len(data))
		if len(data) > (1 << 21) {
			log.Debug().Msgf(" > less than: %vb", 1<<21)
			if err := p.store.Put(ctx, p.currentHash, data); err != nil {
				return err
			}

			//if err := tracker.Set(LObjTrackerPrefix+"/"+hash, []byte(nil)); err != nil {
			//	return err
			//}

			k := path.Join(p.currentHash, "objects", hash[:2], hash[2:])
			if err := p.store.Put(ctx, k, data); err != nil {
				return err
			}

			log.Debug().Str("current_hash", p.currentHash).Msgf("bigNodePatcher")
			p.currentHash = k
			log.Debug().Str("current_hash", p.currentHash).Msgf("bigNodePatcher")
		}

		return nil
	}
}

func (p *Pfg) getRef(ctx context.Context, name string) ([]byte, error) {
	key := path.Join(p.remoteName, name)

	has, err := p.store.Has(ctx, key)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	data, err := p.store.(storage.ReadableStorage).Get(ctx, key)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func isNoLink(err error) bool {
	return strings.Contains(err.Error(), "no link named") || strings.Contains(err.Error(), "no link by that name")
}
