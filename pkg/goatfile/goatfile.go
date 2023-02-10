// Package goatfile provides functionalities to
// unmarshal and parse a Goatfile.
//
// Here you can find the Goatfile specification
// on which basis this parser in built on.
// https://github.com/studio-b12/goat/blob/main/docs/goatfile-spec.md
package goatfile

import (
	"github.com/studio-b12/goat/pkg/util"
)

const (
	sectionNameSetup        = "setup"
	sectionNameSetupEach    = "setup-each"
	sectionNameTests        = "tests"
	sectionNameTeardown     = "teardown"
	sectionNameTeardownEach = "teardown-each"
)

const (
	optionNameQueryParams = "queryparams"
	optionNameHeader      = "header"
	optionNameHeaders     = "headers"
	optionNameBody        = "body"
	optionNameScript      = "script"
	optionNameOptions     = "options"
	abc                   = "asd"
)

// Goatfile holds all sections and
// their requests.
type Goatfile struct {
	Imports []string

	Setup        []Request
	SetupEach    []Request
	Tests        []Request
	Teardown     []Request
	TeardownEach []Request

	Path string
}

// Merge appends all requests in all sections of with
// to the current Goatfile.
func (t *Goatfile) Merge(with Goatfile) {
	t.Setup = append(t.Setup, with.Setup...)
	t.SetupEach = append(t.SetupEach, with.SetupEach...)
	t.Tests = append(t.Tests, with.Tests...)
	t.Teardown = append(t.Teardown, with.Teardown...)
	t.TeardownEach = append(t.TeardownEach, with.TeardownEach...)
}

// String returns the Goatfile as JSON encoded string.
func (t Goatfile) String() string {
	return util.MustJsonMarshalIndent(t)
}

// Opts holds the specific request
// options.
type Opts struct {
	QueryParams map[string]any
	Options     map[string]any
}
