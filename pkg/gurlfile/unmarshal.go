package gurlfile

import "strings"

// Unmarshal takes a raw string of a Gurlfile and tries
// to parse it. Returns the parsed Gurlfile.
func Unmarshal(raw string, currDir string, params ...any) (gf Gurlfile, err error) {

	raw = crlf2lf(raw)

	gf, err = NewParser(strings.NewReader(raw)).Parse()

	return gf, err
}
