package goatfile

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/studio-b12/goat/pkg/errs"
	"github.com/studio-b12/goat/pkg/goatfile/ast"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

// Data provides a getter to receive
// a reader to internal data.
type Data interface {
	// Reader returns a reader to ther internal
	// data or an error. The returned reader
	// might be nil.
	Reader() (io.Reader, error)
}

func DataFromAst(di ast.DataContent, filePath string) (data Data, header http.Header, err error) {
	switch d := di.(type) {
	case ast.NoContent:
		return NoContent{}, nil, nil
	case ast.TextBlock:
		return StringContent(d.Content), nil, nil
	case ast.FileDescriptor:
		fc := FileContent{
			filePath: d.Path,
			currDir:  path.Dir(filePath),
		}
		if d.ContentType != "" {
			header = http.Header{
				"Content-Type": []string{d.ContentType},
			}
		}
		return fc, header, nil
	case ast.RawDescriptor:
		rc := RawContent{
			varName: d.VarName,
		}
		return rc, header, nil
	case ast.FormData:
		boundary, err := randomBoundary()
		if err != nil {
			return nil, nil, err
		}
		header = http.Header{
			"Content-Type": []string{"multipart/form-data; boundary=" + boundary},
		}
		fd := FormData{
			fields:   d.KVList.ToMap(),
			currDir:  path.Dir(filePath),
			boundary: boundary,
		}
		return fd, header, nil
	default:
		return nil, nil, fmt.Errorf("invalid ast data content type: %v", di)
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
	pth, err := joinPath(t.currDir, t.filePath)
	if err != nil {
		return nil, err
	}
	r, err = os.Open(pth)
	return r, err
}

// RawContent can be used for reading byte
// array data
type RawContent struct {
	varName string
	value   any
}

func (t RawContent) Reader() (r io.Reader, err error) {
	rv := reflect.ValueOf(t.value)
	if rv.Kind() == reflect.Slice && rv.Type().Elem().Kind() == reflect.Uint8 {
		r = bytes.NewReader(rv.Bytes())
		return r, nil
	}
	return nil, fmt.Errorf("variable is not a byte array: %v", t.varName)
}

// FormData writes the given key-value pairs into a Multipart Formdata
// encoded reader stream.
type FormData struct {
	fields   map[string]any
	currDir  string
	boundary string
}

func (t FormData) Reader() (io.Reader, error) {
	var b bytes.Buffer

	w := multipart.NewWriter(&b)
	err := w.SetBoundary(t.boundary)
	if err != nil {
		return nil, err
	}
	defer w.Close()

	//goland:noinspection GoDeferInLoop
	for k, v := range t.fields {
		if fd, ok := v.(ast.FileDescriptor); ok {
			filePath, err := joinPath(t.currDir, fd.Path)
			if err != nil {
				return nil, err
			}
			fileName := filepath.Base(filePath)

			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition",
				fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
					quoteEscaper.Replace(k), quoteEscaper.Replace(fileName)))
			if fd.ContentType == "" {
				fd.ContentType = "application/octet-stream"
			}
			h.Set("Content-Type", fd.ContentType)
			fw, err := w.CreatePart(h)
			if err != nil {
				return nil, err
			}

			f, err := os.Open(filePath)
			if err != nil {
				return nil, err
			}
			defer f.Close()
			_, err = io.Copy(fw, f)
			if err != nil {
				return nil, err
			}
			continue
		}

		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Slice && rv.Type().Elem().Kind() == reflect.Uint8 {
			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition",
				fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
					quoteEscaper.Replace(k), quoteEscaper.Replace("binary-data")))
			contentType := http.DetectContentType(rv.Bytes())
			h.Set("Content-Type", contentType)
			fw, err := w.CreatePart(h)
			if err != nil {
				return nil, err
			}
			_, err = fw.Write(rv.Bytes())
			if err != nil {
				return nil, err
			}
			continue
		}

		fw, err := w.CreateFormField(k)
		if err != nil {
			return nil, err
		}
		_, err = fmt.Fprintf(fw, "%v", v)
		if err != nil {
			return nil, err
		}
	}

	return &b, nil
}

func IsNoContent(d Data) bool {
	_, ok := d.(NoContent)
	return ok
}

func joinPath(currDir string, fileName string) (string, error) {
	if strings.HasPrefix(fileName, "~/") {
		currentUser, err := user.Current()
		if err != nil {
			return "", errs.WithPrefix(
				"failed resolving current user for relative path resolution:", err)
		}
		return path.Join(currentUser.HomeDir, fileName[2:]), nil
	}

	if filepath.IsAbs(fileName) {
		return fileName, nil
	}

	return filepath.Join(currDir, fileName), nil
}

func randomBoundary() (string, error) {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", buf[:]), nil
}
