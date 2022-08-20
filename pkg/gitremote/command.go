package gitremote

import "strings"

const (
	CmdList  = "list"
	CmdPush  = "push"
	CmdFetch = "fetch"
)

var DefaultCapabilities = strings.Join([]string{CmdPush, CmdFetch}, "\n")
