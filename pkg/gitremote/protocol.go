package gitremote

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	gitremote "github.com/drgomesp/git-remote-ipldprime"
	ipldgit "github.com/drgomesp/git-remote-ipldprime/core"
	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

type Protocol struct {
	prefix   string
	localDir string

	tracker  *ipldgit.Tracker
	handler  gitremote.ProtocolHandler
	repo     *git.Repository
	lazyWork []func() (string, error)
}

func NewProtocol(prefix string, tracker *ipldgit.Tracker, handler gitremote.ProtocolHandler) (*Protocol, error) {
	log.Info().Msgf("GIT_DIR=%s", os.Getenv("GIT_DIR"))

	localDir, err := GetLocalDir()
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(localDir)
	if err == git.ErrWorktreeNotProvided {
		repoRoot, _ := path.Split(localDir)

		repo, err = git.PlainOpen(repoRoot)
		if err != nil {
			return nil, err
		}
	}

	if err = handler.Initialize(tracker, repo); err != nil {
		return nil, err
	}

	return &Protocol{
		prefix:   prefix,
		handler:  handler,
		repo:     repo,
		localDir: localDir,
		tracker:  tracker,
		lazyWork: make([]func() (string, error), 0),
	}, nil
}

func (p *Protocol) Run(r io.Reader, w io.Writer) (err error) {
	reader := bufio.NewReader(r)
loop:
	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		command = strings.Trim(command, "\n")

		log.Info().Msgf("< %s", command)
		switch {
		case command == "capabilities":
			p.Printf(w, "push\n")
			p.Printf(w, "fetch\n")
			p.Printf(w, "\n")
		case strings.HasPrefix(command, "list"):
			list, err := p.handler.List(strings.HasPrefix(command, "list for-push"))
			if err != nil {
				log.Err(err).Send()
				return err
			}
			for _, e := range list {
				p.Printf(w, "%s\n", e)
			}
			p.Printf(w, "\n")
		case strings.HasPrefix(command, "push "):
			refs := strings.Split(command[5:], ":")
			p.push(refs[0], refs[1], false)
		case strings.HasPrefix(command, "fetch "):
			parts := strings.Split(command, " ")
			if parts[1] != "0000000000000000000000000000000000000000" {
				p.fetch(parts[1], parts[2])
			}
		case command == "":
			fallthrough
		case command == "\n":
			log.Info().Msg("Processing tasks")
			for _, task := range p.lazyWork {
				resp, err := task()
				if err != nil {
					return err
				}
				p.Printf(w, "%s", resp)
			}
			p.Printf(w, "\n")
			p.lazyWork = nil
			break loop
		default:
			return fmt.Errorf("received unknown command %q", command)
		}
	}

	return p.handler.Finish()
}

func (p *Protocol) push(src string, dst string, force bool) {
	_ = force // TODO: handle force push

	p.lazyWork = append(p.lazyWork, func() (string, error) {
		done, err := p.handler.Push(context.Background(), src, dst)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("ok %s\n", done), nil
	})
}

func (p *Protocol) fetch(sha string, ref string) {
	p.lazyWork = append(p.lazyWork, func() (string, error) {
		if sha == "0000000000000000000000000000000000000000" {
			return "", nil
		}

		fetch := p.NewFetch()
		err := fetch.FetchHash(sha)
		if err != nil {
			return "", fmt.Errorf("command fetch: %v", err)
		}

		sha, err := hex.DecodeString(sha)
		if err != nil {
			return "", fmt.Errorf("fetch: %v", err)
		}

		p.tracker.Set(ref, sha)
		return "", nil
	})
}

func (p *Protocol) Printf(w io.Writer, format string, a ...interface{}) {
	if _, err := fmt.Fprintf(w, format, a...); err != nil {
		log.Err(err).Send()
	}
}

func (p *Protocol) NewFetch() *ipldgit.Fetch {
	return ipldgit.NewFetch(p.localDir, p.tracker, p.handler.ProvideBlock)
}
