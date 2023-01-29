package gurlfile

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/studio-b12/gurl/pkg/errs"
	"github.com/studio-b12/gurl/pkg/set"
)

// Unmarshal takes a raw string of a Gurlfile and tries
// to parse it. Returns the parsed Gurlfile.
func Unmarshal(raw string, currDir string) (gf Gurlfile, err error) {
	return unmarshal(os.DirFS("."), raw, currDir, set.Set[string]{})
}

func unmarshal(fSys fs.FS, raw string, currDir string, visited set.Set[string]) (gf Gurlfile, err error) {

	log.Trace().Str("currDir", currDir).Msg("Unmarshalling Gurlfile ...")

	raw = crlf2lf(raw)

	gf, err = NewParser(strings.NewReader(raw)).Parse()
	if err != nil {
		return Gurlfile{}, err
	}

	if !visited.Add(gf.String()) {
		return Gurlfile{}, ErrMultiImport
	}

	var imports Gurlfile
	for _, path := range gf.Imports {
		fullPath := extend(filepath.Join(currDir, path), FileExtension)

		raw, err := fs.ReadFile(fSys, fullPath)
		if err != nil {
			return Gurlfile{}, errs.WithPrefix(
				fmt.Sprintf("failed following import %s:", fullPath), err)
		}

		relativeCurrDir := filepath.Dir(fullPath)
		importGf, err := unmarshal(fSys, string(raw), relativeCurrDir, visited)
		if err != nil {
			return Gurlfile{}, errs.WithPrefix(
				fmt.Sprintf("failed parsing imported file %s:", fullPath), err)
		}

		imports.Merge(importGf)
	}

	imports.Merge(gf)

	return imports, err
}
