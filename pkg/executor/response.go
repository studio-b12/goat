package executor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
			fmt.Errorf("failed reading response body: %s", err.Error())
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
