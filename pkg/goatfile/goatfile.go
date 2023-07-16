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

type SectionName string

const (
	SectionSetup        = SectionName("setup")
	SectionSetupEach    = SectionName("setup-each")
	SectionTests        = SectionName("tests")
	SectionTeardown     = SectionName("teardown")
	SectionTeardownEach = SectionName("teardown-each")
	SectionDefaults     = SectionName("defaults")
)

type optionName string

const (
	optionNameQueryParams = optionName("queryparams")
	optionNameHeader      = optionName("header")
	optionNameHeaders     = optionName("headers")
	optionNameBody        = optionName("body")
	optionNamePreScript   = optionName("prescript")
	optionNameScript      = optionName("script")
	optionNameOptions     = optionName("options")
)

// Goatfile holds all sections and
// their requests.
type Goatfile struct {
	Imports []string

	Defaults *Request

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
	if t.Defaults == nil && with.Defaults != nil {
		t.Defaults = with.Defaults
	} else {
		t.Defaults.Merge(with.Defaults)
	}

	t.Setup = append(t.Setup, with.Setup...)
	t.SetupEach = append(t.SetupEach, with.SetupEach...)
	t.Tests = append(t.Tests, with.Tests...)
	t.Teardown = append(t.Teardown, with.Teardown...)
	t.TeardownEach = append(t.TeardownEach, with.TeardownEach...)

	t.Path = with.Path
}

// String returns the Goatfile as JSON encoded string.
func (t Goatfile) String() string {
	return util.SafeJsonMarshalIndent(t)
}

// Opts holds the specific request
// options.
type Opts struct {
	QueryParams map[string]any
	Options     map[string]any
}
