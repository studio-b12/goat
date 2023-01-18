package gurlfile

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	rxSections         = regexp.MustCompile(`(?m)^#{3,}\s+\w+$`)
	rxRequestSeparator = regexp.MustCompile(`(?m)^-{3,}$`)
	rxHeader           = regexp.MustCompile(`(?m)^([\w\-]+):\s*(.+)$`)
	rxOptionHeader     = regexp.MustCompile(`(?m)^\s*\[([\w\-]+)\]\s*(.*\n)*`)
)

// Unmarshal takes a raw string of a Gurlfile and tries
// to parse it. Returns the parsed Gurlfile.
func Unmarshal(raw string, currDir string, params ...any) (gf Gurlfile, err error) {

	raw = removeComments(raw)

	importPathes, rest, err := parseImports(raw)
	if err != nil {
		return Gurlfile{}, err
	}
	raw = rest

	if currDir == "" {
		currDir, err = os.Getwd()
		if err != nil {
			return Gurlfile{}, fmt.Errorf("failed getting cwd: %s", err.Error())
		}
	}

	for _, path := range importPathes {
		fullPath := extend(filepath.Join(currDir, path), FileExtension)

		raw, err := os.ReadFile(fullPath)
		if err != nil {
			return Gurlfile{}, fmt.Errorf("failed following import %s: %s",
				fullPath, err.Error())
		}

		relativeCurrDir := filepath.Dir(fullPath)
		importGf, err := Unmarshal(string(raw), relativeCurrDir, params...)
		if err != nil {
			return Gurlfile{}, fmt.Errorf("failed parsing imported file %s: %s",
				fullPath, err.Error())
		}

		(&gf).Merge(importGf)
	}

	sections := splitSections(raw)

	for _, section := range sections {
		err = parseSection(section[0], section[1], &gf)
		if err != nil {
			return Gurlfile{}, err
		}
	}

	return gf, nil
}

func splitSections(raw string) [][]string {
	sectionsIndices := rxSections.FindAllStringIndex(raw, -1)

	var sections [][]string

	for i, indices := range sectionsIndices {
		start := indices[0]
		end := len(raw)
		if i+1 < len(sectionsIndices) {
			end = sectionsIndices[i+1][0]
		}

		name := strings.Fields(raw[start:indices[1]])[1]
		content := raw[indices[1]:end]

		sections = append(sections, []string{
			strings.TrimSpace(name),
			strings.TrimSpace(content),
		})
	}

	if len(sections) == 0 {
		sections = append(sections, []string{sectionnameTests, strings.TrimSpace(raw)})
	}

	return sections
}

func parseSection(name, content string, gf *Gurlfile) error {
	requestsRaw := rxRequestSeparator.Split(content, -1)
	requests := make([]Request, 0, len(requestsRaw))

	for i, requestRaw := range requestsRaw {
		ctx := context{section: name, index: i}
		request, err := ctx.parseRequest(requestRaw, nil)
		if err != nil {
			return err
		}
		requests = append(requests, request)
	}

	name = strings.ToLower(name)
	switch name {
	case sectionNameSetup:
		gf.Setup = requests
	case sectionNameSetupEach:
		gf.SetupEach = requests
	case sectionnameTests:
		gf.Tests = requests
	case sectionNameTeardown:
		gf.Teardown = requests
	case sectionNameTeardownEach:
		gf.TeardownEach = requests
	default:
		return newDetailedError(ErrInvalidSectionName,
			"The section name %s is invalid.", name)
	}

	return nil
}

func (t context) parseRequest(requestRaw string, params any) (req Request, err error) {
	defer func() {
		// If an error is returned, wrap the error
		// in a ContextError.
		if err != nil {
			err = t.WrapErr(err)
		}
	}()

	requestRaw = strings.TrimSpace(requestRaw)

	if requestRaw == "" {
		return Request{}, ErrEmptyRequest
	}

	if params != nil {
		var err error
		requestRaw, err = applyTemplate(requestRaw, params)
		if err != nil {
			return Request{}, err
		}
	}

	// Split into sections (separated by one or more empty new lines)
	sectionsSplit := strings.Split(requestRaw, "\n\n")
	sections := make([]string, 0, len(sectionsSplit))
	for _, section := range sectionsSplit {
		section = strings.TrimSpace(section)
		if section != "" {
			sections = append(sections, section)
		}
	}

	if len(sections) == 0 {
		return Request{}, ErrEmptyRequest
	}

	// Part 1: Parse Request Method, URL, Headers and Payload

	lines := strings.Split(sections[0], "\n")

	headerSplit := strings.Fields(lines[0])
	if len(headerSplit) < 2 {
		return Request{}, ErrInvalidHead
	}

	req = newRequest()

	req.context = t
	if params == nil {
		req.raw = requestRaw
	}

	req.Method = headerSplit[0]
	req.URI = headerSplit[1]
	if err != nil {
		return Request{}, fmt.Errorf("invalid URL: %s", err.Error())
	}

	parsingHeaders := true
	bodyBuf := bytes.Buffer{}

	for i := 1; i < len(lines); i++ {
		line := lines[i]

		if parsingHeaders {
			matches := rxHeader.FindAllStringSubmatch(line, -1)
			if len(matches) == 0 {
				parsingHeaders = false
			} else {
				for _, match := range matches {
					req.Header.Set(match[1], match[2])
				}
			}
		}

		if !parsingHeaders {
			// Appending a line break before every line of
			// body content to compensate for the missing
			// line break at the end due to the split.
			err = bodyBuf.WriteByte('\n')
			if err != nil {
				return Request{}, fmt.Errorf("failed appending request body: %s", err.Error())
			}
			_, err := bodyBuf.WriteString(line)
			if err != nil {
				return Request{}, fmt.Errorf("failed appending request body: %s", err.Error())
			}
		}
	}

	if bodyBuf.Len() > 1 {
		// Removing the first line break.
		req.Body = bodyBuf.Bytes()[1:]
	}

	// Part 2: Parse Toml Options

	var optionsB, scriptB strings.Builder
	parsingOptions := true

	for i := 1; i < len(sections); i++ {
		section := sections[i]

		if parsingOptions {
			if !rxOptionHeader.MatchString(section) {
				parsingOptions = false
			} else {
				optionsB.WriteString(section)
			}
		}

		if !parsingOptions {
			scriptB.WriteString(section)
		}
	}

	options, err := parseOptions(optionsB.String())
	if err != nil {
		return Request{}, fmt.Errorf("failed parsing options: %s", err.Error())
	}

	req.Options = options
	req.Script = scriptB.String()

	return req, nil
}

func parseOptions(raw string) (Options, error) {
	raw = strings.ReplaceAll(raw, "{{", "\"{{")
	raw = strings.ReplaceAll(raw, "}}", "}}\"")

	var options Options
	_, err := toml.Decode(raw, &options)
	if err != nil {
		return Options{}, err
	}

	return options, nil
}

func parseImports(raw string) (pathes []string, rest string, err error) {
	lines := strings.Split(raw, "\n")

	for i, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "use") {
			continue
		}

		split := strings.Fields(line)

		var path string
		if len(split) > 1 {
			path = unquote(strings.Join(split[1:], " "))
		}

		if path == "" {
			return nil, "", newDetailedError(ErrInvalidUse,
				"use statement must follow a path to a gurlfile")
		}

		pathes = append(pathes, path)
		lines[i] = ""
	}

	rest = strings.Join(lines, "\n")
	return pathes, rest, err
}
