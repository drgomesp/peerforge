// Package gitremotepfg implements a ProtocolHandler that
// stores git objects as IPLD-structured data in IPFS.
package gitremotepfg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/davecgh/go-spew/spew"
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

type IPFS interface {
	DagPut(data interface{}, inputCodec, storeCodec string) (string, error)
	Add(r io.Reader, options ...interface{}) (string, error)
	PatchLink(root, path, childhash string, create bool) (string, error)
	Get(hash, outdir string) error
	List(path string) ([]ipld.Link, error)
	Cat(path string) (io.ReadCloser, error)
	ResolvePath(path string) (string, error)
}

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
	tracker *core.Tracker

	largeObjs   map[string]string
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
	if p.pushed {
		if err := p.fillMissingLobjs(p.tracker); err != nil {
			return err
		}

		log.Printf("Pushed to IPFS as \x1b[32mipld://%s\x1b[39m\n", p.currentHash)
	}

	return nil
}

func (p *Pfg) ProvideBlock(cid string, tracker *core.Tracker) ([]byte, error) {
	if p.largeObjs == nil {
		if err := p.loadObjectMap(); err != nil {
			return nil, err
		}
	}

	mappedCid, ok := p.largeObjs[cid]
	if !ok {
		return nil, core.ErrNotProvided
	}

	if err := tracker.Set(LObjTrackerPrefix+"/"+cid, []byte(mappedCid)); err != nil {
		return nil, err
	}

	data, err := p.store.Get(context.TODO(), mappedCid)
	if err != nil {
		return nil, errors.New("get error")
	}

	err = p.store.Put(context.TODO(), "raw", data)
	if err != nil {
		return nil, err
	}

	//if realCid != cid {
	//	return nil, fmt.Errorf("unexpected cid for provided block %s != %s", realCid, cid)
	//}

	return data, nil
}

func (p *Pfg) loadObjectMap() error {
	p.largeObjs = map[string]string{}

	links, err := p.store.Get(context.TODO(), path.Join(p.currentHash, HEAD))
	if err != nil {
		if isNoLink(err) {
			return nil
		}

		return err
	}

	spew.Printf("links=%v\n", links)
	//for _, link := range links {
	//	h.largeObjs[link.Name] = link.Hash
	//}

	return nil
}

func (p *Pfg) Capabilities() string {
	return gitremote.DefaultCapabilities
}

func (p *Pfg) List(forPush bool) ([]string, error) {
	ctx := context.Background()
	out := make([]string, 0)

	if !forPush {
		refs, err := p.paths(ctx, nil, p.remoteName, 0)
		if err != nil {
			return nil, err
		}

		for _, ref := range refs {
			switch ref.rType {
			case RefPathHead:
				r := path.Join(strings.Split(ref.path, "/")[1:]...)
				c, err := core.CidFromHex(ref.hash)
				if err != nil {
					return nil, err
				}

				hash, err := core.HexFromCid(c)
				if err != nil {
					return nil, err
				}

				out = append(out, fmt.Sprintf("%s %s", hash, r))
			case RefPathRef:
				r := path.Join(strings.Split(ref.path, "/")[1:]...)
				dest, err := p.getRef(ctx, r)
				if err != nil {
					return nil, err
				}
				out = append(out, fmt.Sprintf("@%s %s", dest, r))
			}

		}
	} else {

	}

	return out, nil
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

			if err := tracker.Set(LObjTrackerPrefix+"/"+hash, []byte(nil)); err != nil {
				return err
			}

			k := path.Join(p.currentHash, "objects", hash[:2], hash[2:])
			if err := p.store.Put(ctx, k, data); err != nil {
				return err
			}

			p.currentHash = k
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

func (p *Pfg) paths(ctx context.Context, api IPFS, pth string, level int) ([]refPath, error) {
	//links, err := api.List(pth)
	//ref, err := p.getRef(ctx, pth)

	//if err != nil {
	//	panic(err)
	//}

	out := make([]refPath, 0)

	//if ref != nil {
	v, err := p.store.Get(ctx, path.Join(pth, HEAD))
	if err != nil {
		return nil, err
	}

	ref, err := p.getRef(ctx, string(v))
	if err != nil {
		return nil, err
	}
	out = append(out, refPath{
		path.Join(pth, string(v)),
		RefPathHead,
		string(ref),
	})

	return out, nil
}
func (p *Pfg) fillMissingLobjs(tracker *core.Tracker) error {
	if p.largeObjs == nil {
		if err := p.loadObjectMap(); err != nil {
			return err
		}
	}

	tracked, err := tracker.ListPrefixed(LObjTrackerPrefix)
	if err != nil {
		return err
	}

	for k, v := range tracked {
		if _, has := p.largeObjs[k]; has {
			continue
		}

		k = strings.TrimPrefix(k, LObjTrackerPrefix+"/")

		p.largeObjs[k] = v
		//p.currentHash, err = p.api.PatchLink(p.currentHash, "objects/"+k, v, true)
		//if err != nil {
		//	return err
		//}
	}

	return nil
}

func isNoLink(err error) bool {
	return strings.Contains(err.Error(), "no link named") || strings.Contains(err.Error(), "no link by that name")
}
