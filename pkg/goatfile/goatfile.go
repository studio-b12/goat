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

type sectionName string

const (
	sectionNameSetup        = sectionName("setup")
	sectionNameSetupEach    = sectionName("setup-each")
	sectionNameTests        = sectionName("tests")
	sectionNameTeardown     = sectionName("teardown")
	sectionNameTeardownEach = sectionName("teardown-each")
)

type optionName string

const (
	optionNameQueryParams = optionName("queryparams")
	optionNameHeader      = optionName("header")
	optionNameHeaders     = optionName("headers")
	optionNameBody        = optionName("body")
	optionNameScript      = optionName("script")
	optionNameOptions     = optionName("options")
	abc                   = optionName("asd")
)

// Goatfile holds all sections and
// their requests.
type Goatfile struct {
	Imports []string

	Setup        []Action
	SetupEach    []Action
	Tests        []Action
	Teardown     []Action
	TeardownEach []Action

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
