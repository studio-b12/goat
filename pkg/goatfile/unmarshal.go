package goatfile

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/studio-b12/goat/pkg/errs"
	"github.com/studio-b12/goat/pkg/set"
)

// Unmarshal takes a raw string of a Goatfile and tries
// to parse it. Returns the parsed Goatfile.
func Unmarshal(raw string, currDir string) (gf Goatfile, err error) {
	currDir = strings.ReplaceAll(currDir, "\\", "/")
	return unmarshal(os.DirFS("."), raw, currDir, set.Set[string]{})
}

func unmarshal(fSys fs.FS, raw string, currDir string, visited set.Set[string]) (gf Goatfile, err error) {

	log.Trace().Str("currDir", currDir).Msg("Unmarshalling Goatfile ...")

	raw = crlf2lf(raw)

	gf, err = NewParser(strings.NewReader(raw), currDir).Parse()
	if err != nil {
		return Goatfile{}, err
	}

	if !visited.Add(gf.String()) {
		return Goatfile{}, ErrMultiImport
	}

	var imports Goatfile
	for _, pth := range gf.Imports {
		fullPath := extend(path.Join(currDir, pth), FileExtension)

		log.Trace().Str("fullPath", fullPath).Msg("Reading import file ...")

		raw, err := fs.ReadFile(fSys, fullPath)
		if err != nil {
			return Goatfile{}, errs.WithPrefix(
				fmt.Sprintf("failed following import %s:", fullPath), err)
		}

		relativeCurrDir := path.Dir(fullPath)
		importGf, err := unmarshal(fSys, string(raw), relativeCurrDir, visited)
		if err != nil {
			return Goatfile{}, errs.WithPrefix(
				fmt.Sprintf("failed parsing imported file %s:", fullPath), err)
		}

		imports.Merge(importGf)
	}

	imports.Merge(gf)

	return imports, err
}
