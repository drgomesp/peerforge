package peerforgecli

import (
	"errors"
)

var (
	ErrRepositoryAlreadyInitialized = errors.New("repository already initialized")
)

const (
	RemoteName     = "peerforge"
	ConfigFileName = ".peerforge.yaml"
	DefaultConfig  = `{"foo": "bar"}`
)
