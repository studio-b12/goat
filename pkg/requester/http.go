package requester

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/zekrotja/rogu/log"
)

var logger = log.Tagged("requester")

// noSetWrapper wraps a cookiejar where no
// cookies can be set by the response.
type noSetWrapper struct {
	http.CookieJar
}

func (t noSetWrapper) SetCookies(u *url.URL, cookies []*http.Cookie) {}

// noGetWrapper wraps a cookiejar where no
// cookies will be sent with the request.
type noGetWrapper struct {
	http.CookieJar
}

func (t noGetWrapper) Cookies(u *url.URL) []*http.Cookie { return nil }

// HttpWithCookies implements Requester with
// the default net/http Client. Also, a map
// of cookie jars is utilized to choose
// between and manipulate the behavior
// of cookie handling.
type HttpWithCookies struct {
	client     *http.Client
	cookieJars map[any]http.CookieJar
}

var _ Requester = (*HttpWithCookies)(nil)

// NewHttpWithCookies returns a new instance of HttpWithCookies.
// cfg is getting passed the instance of http.Client which you
// can configure as you desire.
func NewHttpWithCookies(cfg func(client *http.Client)) *HttpWithCookies {
	var t HttpWithCookies

	t.client = http.DefaultClient

	cfg(t.client)

	t.cookieJars = make(map[any]http.CookieJar)

	return &t
}

func (t HttpWithCookies) Do(req *http.Request, opt Options) (*http.Response, error) {
	jar, err := t.getJar(&opt)
	if err != nil {
		return nil, err
	}

	logger.Trace().Fields(
		"method", req.Method,
		"url", req.URL,
		"header", req.Header,
		"cookies", jar.Cookies(req.URL),
	).Msg("Sending request ...")

	client := *t.client
	client.Jar = jar

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	logger.Trace().Fields(
		"statusCode", res.StatusCode,
		"header", res.Header,
	).Msg("Received response")

	return res, nil
}

// getJar takes a cookiejar from the internal jar map
// by the given key in the options or creates one
// if no jar has already been created.
//
// The returned jar is wrapped by the noSetWrapper
// and/or noGetWrapper depending on the passed
// options.
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
