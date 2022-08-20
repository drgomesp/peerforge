package ipfs

import (
	"io"

	shell "github.com/ipfs/go-ipfs-api"
)

type API interface {
	DagPut(data interface{}, inputCodec, storeCodec string) (string, error)
	Add(r io.Reader, options ...shell.AddOpts) (string, error)
	PatchLink(root, path, childhash string, create bool) (string, error)
	Get(hash, outdir string) error
	List(path string) ([]*shell.LsLink, error)
	Cat(path string) (io.ReadCloser, error)
	ResolvePath(path string) (string, error)
}
