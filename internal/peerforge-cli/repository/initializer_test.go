package repository_test

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/drgomesp/peerforge/internal/peerforge-cli/repository"
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

func TestInitializer_Init(t *testing.T) {
	type args struct {
		dir string
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{name: "init repo (new)", args: args{dir: "a"}, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			abciClientMock := new(abciClientMock)
			abciClientMock.
				On("BroadcastTxCommit", mock.Anything, mock.Anything).
				Return(&coretypes.ResultBroadcastTxCommit{}, nil)

			initializer := repository.NewInitializer(abciClientMock)

			setupTest(t, tt.args.dir)
			err := initializer.Init(tt.args.dir)
			teardownTest(t, tt.args.dir)

			assert.NoError(t, err)
			abciClientMock.AssertExpectations(t)
		})
	}
}

func setupTest(t *testing.T, gitDir string) {
	cwd, _ := os.Getwd()
	os.Setenv("GIT_DIR", path.Join(cwd, gitDir))
}

func teardownTest(t *testing.T, gitDir string) {
	cwd, _ := os.Getwd()
	assert.NoError(t, os.RemoveAll(path.Join(cwd, gitDir)))
}
