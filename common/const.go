package common

import (
	"os"
	"path/filepath"
)

const (
	IF0_VERSION  = "IF0_VERSION"
	ZERO_VERSION = "ZERO_VERSION"
)

var (
	RootPath, _  = os.UserHomeDir()
	If0Dir       = filepath.Join(RootPath, ".if0")
	EnvDir       = filepath.Join(If0Dir, ".environments")
	SnapshotsDir = filepath.Join(If0Dir, ".snapshots")
	If0Default   = filepath.Join(If0Dir, "if0.env")
)

// common flag
var (
	Verbose bool
)