package requester

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// HttpWithCookies implements Requester with
// the default net/http Client. Also, a cookiejar
// is attached to the client to collect and send
// cookies between requests.
type HttpWithCookies struct {
	client     *http.Client
	cookieJars map[any]http.CookieJar
}

type noSetWrapper struct {
	http.CookieJar
}

func (t noSetWrapper) SetCookies(u *url.URL, cookies []*http.Cookie) { return }

type noGetWrapper struct {
	http.CookieJar
}

func (t noGetWrapper) Cookies(u *url.URL) []*http.Cookie { return nil }

var _ Requester = (*HttpWithCookies)(nil)

// NewHttpWithCookies returns a new HttpWithCookies
// instance with the given HTTP client. When no
// client is specified, the http.DefaultClient is
// used.
func NewHttpWithCookies(client ...*http.Client) *HttpWithCookies {
	var t HttpWithCookies

	if len(client) != 0 {
		t.client = client[0]
	} else {
		t.client = http.DefaultClient
	}

	t.cookieJars = make(map[any]http.CookieJar)

	return &t
}

func (t HttpWithCookies) Do(req *http.Request, opt Options) (*http.Response, error) {
	jar, err := t.getJar(&opt)
	if err != nil {
		return nil, err
	}

	client := *t.client
	client.Jar = jar

	return client.Do(req)
}

func (t HttpWithCookies) getJar(opt *Options) (jar http.CookieJar, err error) {
	key := opt.CookieJar
	if key == nil {
		key = "default"
	}

	jar, ok := t.cookieJars[key]
	if !ok {
		jar, err = cookiejar.New(nil)
		if err != nil {
			return nil, err
		}
		t.cookieJars[key] = jar
	}

	if !opt.SendCookies {
		jar = noGetWrapper{jar}
	}
	if !opt.StoreCookies {
		jar = noSetWrapper{jar}
	}

	return jar, nil
}
