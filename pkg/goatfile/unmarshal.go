package goatfile

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/studio-b12/goat/pkg/errs"
	"github.com/studio-b12/goat/pkg/set"
	"github.com/zekrotja/rogu/log"
)

// Unmarshal takes a raw string of a Goatfile and tries
// to parse it. Returns the parsed Goatfile.
func Unmarshal(raw string, fileDir string) (gf Goatfile, err error) {
	fileDir = strings.ReplaceAll(fileDir, "\\", "/")
	return unmarshal(os.DirFS("."), raw, fileDir, set.Set[string]{})
}

func unmarshal(fSys fs.FS, raw string, fileDir string, visited set.Set[string]) (gf Goatfile, err error) {

	log.Trace().Field("fileDir", fileDir).Msg("Unmarshalling Goatfile ...")

	raw = crlf2lf(raw)

	gf, err = NewParser(strings.NewReader(raw), fileDir).Parse()
	if err != nil {
		return Goatfile{}, err
	}
	gf.Path = fileDir

	if !visited.Add(gf.String()) {
		return Goatfile{}, ErrMultiImport
	}

	var imports Goatfile
	for _, pth := range gf.Imports {
		fullPath := Extend(path.Join(path.Dir(fileDir), pth), FileExtension)

		log.Trace().Field("fullPath", fullPath).Msg("Reading import file ...")

		raw, err := fs.ReadFile(fSys, fullPath)
		if err != nil {
			return Goatfile{}, errs.WithPrefix(
				fmt.Sprintf("failed following import %s:", fullPath), err)
		}

		importGf, err := unmarshal(fSys, string(raw), fullPath, visited)
		if err != nil {
			return Goatfile{}, errs.WithPrefix(
				fmt.Sprintf("failed parsing imported file %s:", fullPath), err)
		}

		imports.Merge(importGf)
	}

	imports.Merge(gf)

	return imports, err
}
