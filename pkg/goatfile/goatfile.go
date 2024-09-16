// Package goatfile provides functionalities to
// unmarshal and parse a Goatfile.
//
// Here you can find the Goatfile specification
// on which basis this parser in built on.
// https://github.com/studio-b12/goat/blob/main/docs/goatfile-spec.md
package goatfile

import (
	"errors"
	"github.com/studio-b12/goat/pkg/goatfile/ast"
	"github.com/studio-b12/goat/pkg/util"
	"slices"
)

type SectionName string

const (
	SectionSetup    = SectionName("setup")
	SectionTests    = SectionName("tests")
	SectionTeardown = SectionName("teardown")
	SectionDefaults = SectionName("defaults")
)

type optionName string

const (
	optionNameQueryParams = optionName("queryparams")
	optionNameHeader      = optionName("header")
	optionNameBody        = optionName("body")
	optionNamePreScript   = optionName("prescript")
	optionNameScript      = optionName("script")
	optionNameOptions     = optionName("options")
	optionNameAuth        = optionName("auth")
	optionNameFormData    = optionName("formdata")
)

// Goatfile holds all sections and
// their requests.
type Goatfile struct {
	Imports []string

	Defaults *Request

	Setup    []Action
	Tests    []Action
	Teardown []Action

	Path string
}

func FromAst(astGf *ast.Goatfile) (gf Goatfile, err error) {
	if astGf == nil {
		return gf, errors.New("ast is nil")
	}

	gf.Path = astGf.Dir

	gf.Imports = make([]string, 0, len(astGf.Imports))
	for _, imp := range astGf.Imports {
		gf.Imports = append(gf.Imports, imp.Path)
	}

	if len(astGf.Actions) > 0 {
		gf.Tests = slices.Grow(gf.Tests, len(astGf.Actions))
		for _, act := range astGf.Actions {
			a, err := ActionFromAst(act, astGf.Dir)
			if err != nil {
				return Goatfile{}, err
			}
			gf.Tests = append(gf.Tests, a)
		}
	}

	for _, sect := range astGf.Sections {
		switch s := sect.(type) {
		case ast.SectionDefaults:
			defReq, err := PartialRequestFromAst(s.Request, astGf.Dir)
			if err != nil {
				return Goatfile{}, err
			}
			if gf.Defaults == nil {
				gf.Defaults = defReq
			} else {
				gf.Defaults.Merge(defReq)
			}
		case ast.SectionSetup:
			for _, act := range s.Actions {
				a, err := ActionFromAst(act, astGf.Dir)
				if err != nil {
					return Goatfile{}, err
				}
				gf.Setup = append(gf.Setup, a)
			}
		case ast.SectionTests:
			for _, act := range s.Actions {
				a, err := ActionFromAst(act, astGf.Dir)
				if err != nil {
					return Goatfile{}, err
				}
				gf.Tests = append(gf.Tests, a)
			}
		case ast.SectionTeardown:
			for _, act := range s.Actions {
				a, err := ActionFromAst(act, astGf.Dir)
				if err != nil {
					return Goatfile{}, err
				}
				gf.Teardown = append(gf.Teardown, a)
			}
		}
	}

	return gf, nil
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
	t.Tests = append(t.Tests, with.Tests...)
	t.Teardown = append(t.Teardown, with.Teardown...)

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
	Auth        map[string]any
}
