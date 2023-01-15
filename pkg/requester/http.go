package requester

import (
	"net/http"
	"net/http/cookiejar"
)

// HttpWithCookies implements Requester with
// the default net/http Client. Also, a cookiejar
// is attached to the client to collect and send
// cookies between requests.
type HttpWithCookies struct {
	client *http.Client
}

var _ Requester = (*HttpWithCookies)(nil)

// NewHttpWithCookies returns a new HttpWithCookies
// instance with the given HTTP client. When no
// client is specified, the http.DefaultClient is
// used.
func NewHttpWithCookies(client ...*http.Client) (*HttpWithCookies, error) {
	var t HttpWithCookies

	if len(client) != 0 {
		t.client = client[0]
	} else {
		t.client = http.DefaultClient
	}

	var err error
	t.client.Jar, err = cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (t HttpWithCookies) Do(req *http.Request) (*http.Response, error) {
	return t.client.Do(req)
}
