package gitremote

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	gitremote "github.com/peerforge/git-remote-ipldprime"
	ipldgit "github.com/peerforge/git-remote-ipldprime/core"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	getwd, err := os.Getwd()
	if err != nil {
		return
	}

	_ = os.Setenv("GIT_DIR", getwd)
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

var _ gitremote.ProtocolHandler = &handlerMock{}

type handlerMock struct {
	mock.Mock
}

func (h *handlerMock) ProvideBlock(identifier string, tracker *ipldgit.Tracker) ([]byte, error) {
	//TODO implement me
	return nil, nil
}

func (h *handlerMock) Initialize(tracker *ipldgit.Tracker, repo *git.Repository) error {
	return nil
}

func (h *handlerMock) Finish() error {
	return nil
}

func (h *handlerMock) Capabilities() string {
	return DefaultCapabilities
}

func (h *handlerMock) List(forPush bool) ([]string, error) {
	args := h.Called(forPush)
	return args.Get(0).([]string), args.Error(1)
}

func (h *handlerMock) Push(ctx context.Context, localRef string, remoteRef string) (string, error) {
	args := h.Called(localRef, remoteRef)
	return args.String(0), args.Error(1)
}

func (h *handlerMock) Fetch(sha, ref string) error {
	args := h.Called(sha, ref)
	return args.Error(0)
}

func Test_Protocol(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
		err  error
		mock func(m *handlerMock)
	}{
		{
			name: "empty",
			in:   "",
		},
		{
			name: "capabilities",
			in:   "capabilities",
			out:  DefaultCapabilities,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerMock := new(handlerMock)
			if tt.mock != nil {
				tt.mock(handlerMock)
			}

			proto := &Protocol{
				prefix:  "origin",
				handler: handlerMock,
			}

			reader := strings.NewReader(tt.in + "\n")
			var writer bytes.Buffer
			if err := proto.Run(reader, &writer); err != nil {
				if tt.err != io.EOF && tt.err != nil {
					assert.Equal(t, tt.err, err)
				}
			}

			handlerMock.AssertExpectations(t)

			got := strings.TrimSpace(writer.String())

			assert.Equal(t, tt.out, got)
		})
	}
}
