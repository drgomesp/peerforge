package gitremotepfg_test

import (
	"os"
	"path"
	"strings"
	"testing"

	gitremotego "github.com/drgomesp/git-remote-go"
	gitremotepfg "github.com/drgomesp/peerforge/internal/git-remote-pfg"
	"github.com/go-git/go-git/v5"
	"github.com/ipfs-shipyard/git-remote-ipld/core"
	ipfs "github.com/ipfs/go-ipfs-api"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const EmptyRepo = "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

}

func TestIPFS_Capabilities(t *testing.T) {
	type fields struct {
		remoteName string
	}

	var tests = []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test capabilities",
			want: gitremotego.DefaultCapabilities,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remote := setupTest(t, EmptyRepo)
			capab := remote.Capabilities()

			log.Trace().Msgf(">\n%s", capab)
			assert.Equal(t, tt.want, capab)
		})
	}
}

func TestIPFS_Push(t *testing.T) {
	type fields struct {
		localRef, remoteRef string
	}

	var tests = []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "push", fields: fields{
			localRef:  "refs/heads/main",
			remoteRef: "refs/heads/main",
		}, want: `refs/heads/main`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remote := setupTest(t, EmptyRepo)

			push, err := remote.Push(
				tt.fields.localRef,
				tt.fields.remoteRef,
			)
			assert.NoError(t, remote.Finish())
			log.Trace().Msgf(">\n%s", push)

			assert.NoError(t, err)
			assert.NotEmpty(t, push)

			assert.Equal(t, tt.want, push)
		})
	}
}

func TestIPFS_List(t *testing.T) {
	type fields struct {
		remoteName string
	}

	var tests = []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "list",
			fields: fields{
				remoteName: "QmRwnExWq7ypuj6wh8uKHbV5VRpLc73hoGPQiWzYS7r3j8",
			},
			want: []string{
				"refs/heads/main HEAD",
				"ada5ec06cbbdc9616a6e4a7cd43a7b078936368e refs/heads/main",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remote := setupTest(t, tt.fields.remoteName)

			list, err := remote.List(false)
			assert.NoError(t, remote.Finish())
			log.Trace().Msgf(">\n%s", strings.Join(list, "\n"))

			assert.NoError(t, err)
			assert.NotEmpty(t, list)
			assert.Equal(t, tt.want, list)
		})
	}
}

func TestIPFS_Fetch(t *testing.T) {
	type fields struct {
		remoteName string
	}

	var tests = []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "fetch ada5ec06cbbdc9616a6e4a7cd43a7b078936368e refs/heads/main",
			fields: fields{
				remoteName: "QmRwnExWq7ypuj6wh8uKHbV5VRpLc73hoGPQiWzYS7r3j8",
			},
			want: []string{
				"refs/heads/main HEAD",
				"ada5ec06cbbdc9616a6e4a7cd43a7b078936368e refs/heads/main",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remote := setupTest(t, tt.fields.remoteName)

			list, err := remote.List(false)
			assert.NoError(t, remote.Finish())
			log.Trace().Msgf(">\n%s", strings.Join(list, "\n"))

			assert.NoError(t, err)
			assert.NotEmpty(t, list)
			assert.Equal(t, tt.want, list)
		})
	}
}

func setupTest(t *testing.T, remoteName string) *gitremotepfg.PeerforgeRemote {
	cwd, err := os.Getwd()
	if err != nil {
		log.Err(err).Send()
	}

	_ = os.Setenv("GIT_DIR", path.Join(cwd, "..", "..", "testutil", "git"))
	_ = os.Setenv(ipfs.EnvDir, "localhost:5001")

	localDir, err := core.GetLocalDir()
	if err != nil {
		t.Fatal(err)
	}

	repo, err := git.PlainOpen(localDir)
	if err == git.ErrWorktreeNotProvided {
		repoRoot, _ := path.Split(localDir)

		repo, err = git.PlainOpen(repoRoot)
		if err != nil {
			t.Fatal(err)
		}
	}

	tracker, err := core.NewTracker(localDir)
	if err != nil {
		t.Fatal(err)
	}

	remote, err := gitremotepfg.NewPeerforgeRemote(remoteName)
	assert.NoError(t, err)
	assert.NotNil(t, remote)

	assert.NoError(t, remote.Initialize(tracker, repo))

	return remote
}
