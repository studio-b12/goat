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
	Header        http.Header
	ContentLength int64
	Body          string
	BodyJson      any
}

// FromHttpResponse builds a Response from the
// given Http Response reference.
func FromHttpResponse(resp *http.Response) (Response, error) {
	var r Response

	r.StatusCode = resp.StatusCode
	r.Status = resp.Status
	r.Proto = resp.Proto
	r.ProtoMajor = resp.ProtoMajor
	r.ProtoMinor = resp.ProtoMinor
	r.Header = resp.Header
	r.ContentLength = resp.ContentLength

	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{},
			errs.WithPrefix("failed reading response body:", err)
	}

	if len(d) > 0 {
		r.Body = string(d)

		var bodyJson any
		err = json.Unmarshal(d, &bodyJson)
		if err == nil {
			r.BodyJson = bodyJson
		}
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

	if t.BodyJson != nil {
		fmt.Fprintln(&sb)
		enc := json.NewEncoder(&sb)
		enc.SetIndent("", "  ")
		// This shouldn't error because it was decoded by
		// via json.Unmarshal before.
		enc.Encode(t.BodyJson)
	} else if len(t.Body) != 0 {
		fmt.Fprintln(&sb)
		sb.WriteString(t.Body)
	}

	return sb.String()
}
