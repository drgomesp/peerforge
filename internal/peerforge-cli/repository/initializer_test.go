package repository_test

import (
	"os"
	"path"
	"testing"

	"github.com/drgomesp/peerforge/internal/peerforge-cli/repository"
	"github.com/stretchr/testify/assert"
)

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
			initializer := repository.NewInitializer(nil)

			setupTest(t, tt.args.dir)
			err := initializer.Init(tt.args.dir)
			teardownTest(t, tt.args.dir)

			assert.NoError(t, err)
		})
	}
}

func setupTest(t *testing.T, gitDir string) {
	cwd, _ := os.Getwd()
	assert.NoError(t, os.Setenv("GIT_DIR", path.Join(cwd, gitDir)))
}

func teardownTest(t *testing.T, gitDir string) {
	cwd, _ := os.Getwd()
	assert.NoError(t, os.RemoveAll(path.Join(cwd, gitDir)))
}
