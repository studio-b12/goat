package goatfile

import (
	"bytes"
	"fmt"
	"github.com/studio-b12/goat/pkg/goatfile/ast"
	"io"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/studio-b12/goat/pkg/errs"
)

// Data provides a getter to receive
// a reader to internal data.
type Data interface {
	// Reader returns a reader to ther internal
	// data or an error. The returned reader
	// might be nil.
	Reader() (io.Reader, error)
}

func DataFromAst(di ast.DataContent, filePath string) (Data, error) {
	switch d := di.(type) {
	case ast.NoContent:
		return NoContent{}, nil
	case ast.TextBlock:
		return StringContent(d.Content), nil
	case ast.FileDescriptor:
		return FileContent{
			filePath: d.Path,
			currDir:  path.Dir(filePath),
		}, nil
	default:
		return nil, fmt.Errorf("invalid ast data content type: %v", di)
	}
}

// NoContent implements Data containing no content.
type NoContent struct{}

func (t NoContent) Reader() (io.Reader, error) {
	return nil, nil
}

// StringContent stores data as a string.
type StringContent string

func (t StringContent) Reader() (io.Reader, error) {
	return bytes.NewBufferString(string(t)), nil
}

// FileContent provides a getter which opens
// a file with the stored path and returns
// it as reader.
type FileContent struct {
	filePath string
	currDir  string
}

func (t FileContent) Reader() (r io.Reader, err error) {
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

func IsNoContent(d Data) bool {
	_, ok := d.(NoContent)
	return ok
}
