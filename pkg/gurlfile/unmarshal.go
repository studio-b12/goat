package gurlfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Unmarshal takes a raw string of a Gurlfile and tries
// to parse it. Returns the parsed Gurlfile.
func Unmarshal(raw string, currDir string) (gf Gurlfile, err error) {

	raw = crlf2lf(raw)

	gf, err = NewParser(strings.NewReader(raw)).Parse()

	var imports Gurlfile
	for _, path := range gf.Imports {
		fullPath := extend(filepath.Join(currDir, path), FileExtension)

		raw, err := os.ReadFile(fullPath)
		if err != nil {
			return Gurlfile{}, fmt.Errorf("failed following import %s: %s",
				fullPath, err.Error())
		}

		relativeCurrDir := filepath.Dir(fullPath)
		importGf, err := Unmarshal(string(raw), relativeCurrDir)
		if err != nil {
			return Gurlfile{}, fmt.Errorf("failed parsing imported file %s: %s",
				fullPath, err.Error())
		}

		(&imports).Merge(importGf)
	}

	(&imports).Merge(gf)

	return imports, err
}
