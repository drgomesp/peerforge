package gitremotepfg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/drgomesp/git-remote-go"
	peerforge "github.com/drgomesp/peerforge/pkg"
	"github.com/drgomesp/peerforge/pkg/event"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/uuid"
	"github.com/ipfs-shipyard/git-remote-ipld/core"
	"github.com/ipfs/go-cid"
	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/rs/zerolog/log"
	"github.com/tendermint/tendermint/rpc/client"
	gitv4 "gopkg.in/src-d/go-git.v4"
)

const (
	LargeObjectDir    = "objects"
	LobjTrackerPrefix = "//lobj"
)

const (
	RefPathHead = iota
	RefPathRef
)

type refPath struct {
	path  string
	rType int

	hash string
}

var _ gitremotego.ProtocolHandler = &PeerforgeRemote{}

type PeerforgeRemote struct {
	ipfs *ipfs.Shell
	repo *git.Repository
	abci client.ABCIClient

	tracker                 *core.Tracker
	didPush                 bool
	largeObjs               map[string]string
	remoteName, currentHash string
	localDir                string
}

func NewPeerforgeRemote(abci client.ABCIClient, ipfsPath string, remoteName string) (*PeerforgeRemote, error) {
	ipfsShell := ipfs.NewShell(ipfsPath)

	if ipfsShell == nil {
		return nil, errors.New("failed to initialize protocol shell")
	}

	cwd, _ := os.Getwd()

	localDir, err := gitremotego.GetLocalDir()
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

	return &PeerforgeRemote{ipfs: ipfsShell, repo: repo, abci: abci, remoteName: remoteName}, nil
}

func (p *PeerforgeRemote) Initialize(tracker *core.Tracker, repo *git.Repository) error {
	p.repo = repo
	p.currentHash = p.remoteName

	localDir, err := gitremotego.GetLocalDir()
	if err != nil {
		return err
	}

	repo, err = git.PlainOpen(localDir)
	if err != nil {
		return err
	}

	p.localDir = localDir
	p.repo = repo
	p.tracker = tracker

	return nil
}
func (p *PeerforgeRemote) Finish() error {
	//TODO: publish
	if p.didPush {
		if err := p.fillMissingLobjs(p.tracker); err != nil {
			return err
		}

		log.Info().Msgf("Pushed to pfg://%s", p.currentHash)
	}

	return nil
}

func (p *PeerforgeRemote) fillMissingLobjs(tracker *core.Tracker) error {
	if p.largeObjs == nil {
		if err := p.loadObjectMap(); err != nil {
			return err
		}
	}

	tracked, err := tracker.ListPrefixed(LobjTrackerPrefix)
	if err != nil {
		return err
	}

	for k, v := range tracked {
		log.Debug().Msgf("tracked= %v=%v", k, v)

		if _, has := p.largeObjs[k]; has {
			continue
		}

		k = strings.TrimPrefix(k, LobjTrackerPrefix+"/")

		p.largeObjs[k] = v
		p.currentHash, err = p.ipfs.PatchLink(p.currentHash, "objects/"+k, v, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PeerforgeRemote) loadObjectMap() error {
	p.largeObjs = map[string]string{}

	links, err := p.ipfs.List(p.currentHash + "/" + LargeObjectDir)
	if err != nil {
		//TODO: Find a better way with coreapi
		if isNoLink(err) {
			return nil
		}
		return err
	}

	for _, link := range links {
		p.largeObjs[link.Name] = link.Hash
	}

	return nil
}

func (p *PeerforgeRemote) ProvideBlock(identifier string, tracker *core.Tracker) ([]byte, error) {
	if p.largeObjs == nil {
		if err := p.loadObjectMap(); err != nil {
			return nil, err
		}
	}

	mappedCid, ok := p.largeObjs[identifier]
	if !ok {
		return nil, core.ErrNotProvided
	}

	if err := p.tracker.Set(LobjTrackerPrefix+"/"+identifier, []byte(mappedCid)); err != nil {
		return nil, err
	}

	r, err := p.ipfs.Cat(fmt.Sprintf("/ipfs/%s", mappedCid))
	if err != nil {
		return nil, errors.New("cat error")
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	realCid, err := p.ipfs.DagPut(data, "raw", "git")
	if err != nil {
		return nil, err
	}

	if realCid != identifier {
		return nil, fmt.Errorf("unexpected cid for provided block %s != %s", realCid, identifier)
	}

	return data, nil
}

func (p *PeerforgeRemote) Capabilities() string {
	return gitremotego.DefaultCapabilities
}

func (p *PeerforgeRemote) List(forPush bool) ([]string, error) {
	out := make([]string, 0)

	if !forPush {
		refs, err := p.paths(p.ipfs, p.remoteName, 0)
		if err != nil {
			return nil, err
		}

		for _, ref := range refs {
			switch ref.rType {
			case RefPathHead:
				r := path.Join(strings.Split(ref.path, "/")[1:]...)
				c, err := cid.Parse(ref.hash)
				if err != nil {
					return nil, err
				}

				hash, err := gitremotego.HexFromCid(c)
				if err != nil {
					return nil, err
				}

				out = append(out, fmt.Sprintf("%s %s", hash, r))
			case RefPathRef:
				r := path.Join(strings.Split(ref.path, "/")[1:]...)
				dest, err := p.getRef(r)
				if err != nil {
					return nil, err
				}

				out = append(out, fmt.Sprintf("%s %s", dest, r))
			}
		}
	} else {
		it, err := p.repo.Branches()
		if err != nil {
			return nil, err
		}

		err = it.ForEach(func(ref *plumbing.Reference) error {
			remoteRef := "0000000000000000000000000000000000000000"

			localRef, err := p.ipfs.ResolvePath(path.Join(p.currentHash, ref.Name().String()))
			if err != nil && !isNoLink(err) {
				return err
			}
			if err == nil {
				refCid, err := cid.Parse(localRef)
				if err != nil {
					return err
				}

				remoteRef, err = gitremotego.HexFromCid(refCid)
				if err != nil {
					return err
				}
			}

			out = append(out, fmt.Sprintf("%s %s", remoteRef, ref.Name()))

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func (p *PeerforgeRemote) Push(local string, remote string) (string, error) {
	p.didPush = true

	localRef, err := p.repo.Reference(plumbing.ReferenceName(local), true)
	if err != nil {
		return "", fmt.Errorf("command push: %v", err)
	}

	headHash := localRef.Hash().String()
	repo, err := gitv4.PlainOpen(p.localDir)
	if err == git.ErrWorktreeNotProvided {
		repoRoot, _ := path.Split(p.localDir)

		repo, err = gitv4.PlainOpen(repoRoot)
		if err != nil {
			return "", err
		}
	}

	log.Debug().Msgf("remote=%s ref=%s", remote, headHash)
	push := core.NewPush(p.localDir, p.tracker, repo)
	push.NewNode = p.bigNodePatcher(p.tracker)

	log.Debug().Msgf("PushHash(%s)", headHash)
	err = push.PushHash(headHash)
	if err != nil {
		return "", fmt.Errorf("push: %v", err)
	}

	hash := localRef.Hash()

	log.Debug().Msgf("tracking %v", hash.String())
	p.tracker.Set(remote, (&hash)[:])

	c, err := core.CidFromHex(headHash)
	if err != nil {
		return "", fmt.Errorf("push: %v", err)
	}

	//patch object
	p.currentHash, err = p.ipfs.PatchLink(p.currentHash, remote, c.String(), true)
	if err != nil {
		return "", fmt.Errorf("push: %v", err)
	}

	head, err := p.getRef("HEAD")
	if err != nil {
		return "", fmt.Errorf("push: %v", err)
	}

	if head == "" {
		log.Debug().Msgf("push(%v, %v)", local, remote)
		headRef, err := p.ipfs.Add(strings.NewReader("refs/heads/main"))
		if err != nil {
			return "", fmt.Errorf("push: %v", err)
		}

		p.currentHash, err = p.ipfs.PatchLink(p.currentHash, "HEAD", headRef, true)
		if err != nil {
			return "", fmt.Errorf("push: %v", err)
		}
	}

	log.Debug().Msgf("sending event transaction...")
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

	res, err := p.abci.BroadcastTxCommit(context.Background(), data)
	if res.CheckTx.IsErr() || res.DeliverTx.IsErr() {
		log.Debug().Msgf(err.Error())
		log.Err(err).Send()
		return "", err
	}

	return local, nil
}

// bigNodePatcher returns a function which patches large object mapping into
// the resulting object
func (p *PeerforgeRemote) bigNodePatcher(tracker *core.Tracker) func(cid.Cid, []byte) error {
	return func(hash cid.Cid, data []byte) error {
		if len(data) > (1 << 21) {
			c, err := p.ipfs.Add(bytes.NewReader(data))
			if err != nil {
				return err
			}

			if err := tracker.Set(LobjTrackerPrefix+"/"+hash.String(), []byte(c)); err != nil {
				return err
			}

			log.Debug().Msgf("bigNodePatcher PatchLink(%s)", hash.String())

			p.currentHash, err = p.ipfs.PatchLink(p.currentHash, "objects/"+hash.String(), c, true)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func (p *PeerforgeRemote) getRef(name string) (string, error) {
	r, err := p.ipfs.Cat(path.Join(p.remoteName, name))
	if err != nil {
		if isNoLink(err) {
			return "", nil
		}
		return "", err
	}
	defer r.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(r)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (p *PeerforgeRemote) paths(api *ipfs.Shell, pathStr string, level int) ([]refPath, error) {
	links, err := api.List(pathStr)
	if err != nil {
		return nil, err
	}

	out := make([]refPath, 0)
	for _, link := range links {
		switch link.Type {
		case ipfs.TDirectory:
			if level == 0 && link.Name == LargeObjectDir {
				continue
			}

			sub, err := p.paths(api, path.Join(pathStr, link.Name), level+1)
			if err != nil {
				return nil, err
			}
			out = append(out, sub...)
		case ipfs.TFile:
			out = append(out, refPath{path.Join(pathStr, link.Name), RefPathRef, link.Hash})
		case -1, 0: //unknown, assume git node
			out = append(out, refPath{path.Join(pathStr, link.Name), RefPathHead, link.Hash})
		default:
			return nil, fmt.Errorf("unexpected link type %d", link.Type)
		}
	}
	return out, nil
}

func isNoLink(err error) bool {
	return strings.Contains(err.Error(), "no link named") || strings.Contains(err.Error(), "no link by that name")
}
