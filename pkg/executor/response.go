package executor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/studio-b12/goat/pkg/errs"
)

// Response is the model passed into the engine
// state containing the requests repsonse data.
type Response struct {
	StatusCode    int
	Status        string
	Proto         string
	ProtoMajor    int
	ProtoMinor    int
	Header        map[string][]string
	ContentLength int64
	BodyRaw       []byte
	Body          any
}

// FromHttpResponse builds a Response from the
// given Http Response reference.
func FromHttpResponse(resp *http.Response, options map[string]any) (Response, error) {
	var r Response

	r.StatusCode = resp.StatusCode
	r.Status = resp.Status
	r.Proto = resp.Proto
	r.ProtoMajor = resp.ProtoMajor
	r.ProtoMinor = resp.ProtoMinor
	r.Header = resp.Header
	r.ContentLength = resp.ContentLength

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{},
			errs.WithPrefix("failed reading response body:", err)
	}

	// Try to parse body depending on 'repsponsetype' option.
	// If 'responsetype' is not set, try to use response
	// Content-Type header instead.
	if len(data) > 0 {
		r.BodyRaw = data
		responseType, ok := options["responsetype"].(string)
		if !ok {
			contentTypeHeader, ok := r.Header["Content-Type"]
			if ok {
				responseType = contentTypeHeader[0]
			}
		}

		// responseType 'raw' prevents body parsing
		// if required and assigns raw bytes to 'Body'
		if responseType == "raw" {
			r.Body = r.BodyRaw
			return r, nil
		}

		parsedBody, err := parseBody(data, responseType)
		if err != nil {
			return Response{}, errs.WithPrefix("failed parsing body:", err)
		}
		r.Body = parsedBody
	}

	return r, nil
}

func (t Response) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s\n", t.Status)
	for key, vals := range t.Header {
		for _, val := range vals {
			fmt.Fprintf(&sb, "%s: %s\n", key, val)
		}
	}

	if t.Body != nil && json.Valid(t.BodyRaw) {
		enc := json.NewEncoder(&sb)
		enc.SetIndent("", "  ")
		// This shouldn't error because it was decoded by
		// via json.Unmarshal before.
		enc.Encode(t.Body)
	} else if len(t.BodyRaw) != 0 {
		sb.WriteString(string(t.BodyRaw))
	}

	return sb.String()
}
