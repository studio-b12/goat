package goatfile

import (
	"bytes"
	"io"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/studio-b12/goat/pkg/errs"
)

type Data interface {
	Reader() (io.Reader, error)
}

type NoData struct{}

func (t NoData) Reader() (io.Reader, error) {
	return nil, nil
}

type StringData string

func (t StringData) Reader() (io.Reader, error) {
	return bytes.NewBufferString(string(t)), nil
}

type FileData struct {
	filePath string
	currDir  string
}

func (t FileData) Reader() (r io.Reader, err error) {
	pth := t.filePath

	if strings.HasPrefix(pth, "~/") {
		currentUser, err := user.Current()
		if err != nil {
			return nil, errs.WithPrefix(
				"failed resolving current user for relative path resolution:", err)
		}
		pth = path.Join(currentUser.HomeDir, pth[2:])
	} else if !filepath.IsAbs(pth) {
		pth = path.Join(t.currDir, pth)
	}

	r, err = os.Open(pth)
	return r, err
}
