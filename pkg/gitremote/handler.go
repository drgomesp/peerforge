package gitremote

import (
	ipldgit "github.com/drgomesp/git-remote-ipld/core"
	"github.com/go-git/go-git/v5"
)

type ProtocolHandler interface {
	Initialize(tracker *ipldgit.Tracker, repo *git.Repository) error
	Capabilities() string
	List(forPush bool) ([]string, error)
	Push(localRef string, remoteRef string) (string, error)
	ProvideBlock(identifier string, tracker *ipldgit.Tracker) ([]byte, error)
	Finish() error
}
