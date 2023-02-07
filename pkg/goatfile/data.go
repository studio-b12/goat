package goatfile

import (
	"bytes"
	"io"
	"os"
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

type FileData string

func (t FileData) Reader() (r io.Reader, err error) {
	r, err = os.Open(string(t))
	return r, err
}
