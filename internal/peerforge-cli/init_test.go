package peerforgecli_test

import (
	"os"
	"path"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	peerforgecli "github.com/drgomesp/peerforge/internal/peerforge-cli"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func setupTest(t *testing.T, gitDir string) {
	cwd, _ := os.Getwd()
	os.Setenv("GIT_DIR", path.Join(cwd, "..", "..", gitDir))
}

func teardownTest(t *testing.T, gitDir string) {
	cwd, _ := os.Getwd()
	assert.NoError(t, os.RemoveAll(path.Join(cwd, gitDir)))
}

func TestInitRepository(t *testing.T) {
	type args struct {
		dir string
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{name: "init repo (create dir)", args: args{dir: "foo"}, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, peerforgecli.Init(tt.args.dir))
			teardownTest(t, tt.args.dir)
		})
	}
}
