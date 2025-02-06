package util

import (
	"io/fs"
	"os"
)

// RootFs implements fs.FS for the local file system. It basically
// aliases (&RootFs{}).Open(path) to os.Open(path).
type RootFs struct{}

var _ fs.FS = (*RootFs)(nil)

func (t *RootFs) Open(name string) (fs.File, error) {
	return os.Open(name)
}
