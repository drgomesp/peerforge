package gitremotepfg_test

import (
	"context"
	"os"
	"path"
	"strings"
	"testing"

	ipldgit "github.com/drgomesp/git-remote-ipld/core"
	gitremotepfg "github.com/drgomesp/peerforge/internal/git-remote-pfg"
	"github.com/drgomesp/peerforge/pkg/gitremote"
	"github.com/go-git/go-git/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/rpc/coretypes"
	"github.com/tendermint/tendermint/types"
)

type abciClientMock struct {
	mock.Mock
}

func (a abciClientMock) ABCIInfo(ctx context.Context) (*coretypes.ResultABCIInfo, error) {
	args := a.Called(ctx)
	return args.Get(0).(*coretypes.ResultABCIInfo), args.Error(1)
}

func (a abciClientMock) ABCIQuery(ctx context.Context, path string, data bytes.HexBytes) (*coretypes.ResultABCIQuery, error) {
	args := a.Called(ctx, path, data)
	return args.Get(0).(*coretypes.ResultABCIQuery), args.Error(1)
}

func (a abciClientMock) ABCIQueryWithOptions(ctx context.Context, path string, data bytes.HexBytes, opts client.ABCIQueryOptions) (*coretypes.ResultABCIQuery, error) {
	args := a.Called(ctx, path, data)
	return args.Get(0).(*coretypes.ResultABCIQuery), args.Error(1)
}

func (a abciClientMock) BroadcastTxCommit(ctx context.Context, tx types.Tx) (*coretypes.ResultBroadcastTxCommit, error) {
	args := a.Called(ctx, tx)
	return args.Get(0).(*coretypes.ResultBroadcastTxCommit), args.Error(1)
}

func (a abciClientMock) BroadcastTxAsync(ctx context.Context, tx types.Tx) (*coretypes.ResultBroadcastTx, error) {
	args := a.Called(ctx, tx)
	return args.Get(0).(*coretypes.ResultBroadcastTx), args.Error(1)
}

func (a abciClientMock) BroadcastTxSync(ctx context.Context, tx types.Tx) (*coretypes.ResultBroadcastTx, error) {
	args := a.Called(ctx, tx)
	return args.Get(0).(*coretypes.ResultBroadcastTx), args.Error(1)
}

const EmptyRepo = "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cwd, err := os.Getwd()
	if err != nil {
		log.Err(err).Send()
	}
	os.Setenv("GIT_DIR", path.Join(cwd, "testutil", "git"))
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
			want: gitremote.DefaultCapabilities,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			abciClientMock := new(abciClientMock)
			remote := setupTest(t, abciClientMock, EmptyRepo)
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
			abciClientMock := new(abciClientMock)
			remote := setupTest(t, abciClientMock, EmptyRepo)

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
			abciClientMock := new(abciClientMock)
			remote := setupTest(t, abciClientMock, tt.fields.remoteName)

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
			abciClientMock := new(abciClientMock)
			remote := setupTest(t, abciClientMock, tt.fields.remoteName)

			list, err := remote.List(false)
			assert.NoError(t, remote.Finish())
			log.Trace().Msgf(">\n%s", strings.Join(list, "\n"))
			teardownTest(t)

			assert.NoError(t, err)
			assert.NotEmpty(t, list)
			assert.Equal(t, tt.want, list)
		})
	}
}

func setupTest(t *testing.T, abciClient client.ABCIClient, remoteName string) *gitremotepfg.Pfg {
	cwd, err := os.Getwd()
	if err != nil {
		log.Err(err).Send()
	}

	_ = os.Setenv("GIT_DIR", path.Join(cwd, "..", "..", "..", "testutil", "git"))

	localDir, err := ipldgit.GetLocalDir()
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

	tracker, err := ipldgit.NewTracker(localDir)
	if err != nil {
		t.Fatal(err)
	}

	ipfsPath := "localhost:45005" // brave browser IPFS
	remote, err := gitremotepfg.NewPfg(abciClient, ipfsPath, remoteName)
	assert.NoError(t, err)
	assert.NotNil(t, remote)

	assert.NoError(t, remote.Initialize(tracker, repo))

	return remote
}

func teardownTest(t *testing.T) {
	cwd, _ := os.Getwd()
	gitDir := path.Join(cwd, "..", "..", "testutil", "gitremote")
	assert.NoError(t, os.RemoveAll(path.Join(cwd, gitDir)))
}
