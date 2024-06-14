package goatfile

import (
	"testing"
)

func FuzzParseUse(f *testing.F) {
	f.Add(" foo.goat")
	f.Add(" foo")
	f.Add(" foo/bar/baz.goat")
	f.Add(" foo/bar/baz")
	f.Fuzz(func(t *testing.T, s string) {
		p := stringParser(s)
		_, _ = p.parseUse()
	})
}

func FuzzParseExecute(f *testing.F) {
	f.Add(` foo.goat (hello="foo") return (foo as bar)`)
	f.Add(` foo.goat (foo="foo" bar="bar") return (foo as bar)`)
	f.Add(` foo.goat (
	foo="foo" 
	bar="bar"
) return (
	foo as bar
	baz as fuz
)`)
	f.Fuzz(func(t *testing.T, s string) {
		p := stringParser(s)
		_, _, _ = p.parseExecute()
	})
}
