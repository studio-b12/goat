package requester

import "net/http"

type Options struct {
	CookieJar    any
	StoreCookies bool
	SendCookies  bool
}

func OptionsFromMap(m map[string]any) Options {
	opt := Options{
		CookieJar:    "default",
		StoreCookies: true,
		SendCookies:  true,
	}

	if v, ok := m["cookiejar"]; ok {
		opt.CookieJar = v
	}
	if v, ok := m["storecookies"].(bool); ok {
		opt.StoreCookies = v
	}
	if v, ok := m["sendcookies"].(bool); ok {
		opt.SendCookies = v
	}

	return opt
}

// Requester defines a service to perform HTTP
// requests.
type Requester interface {
	// Do takes a HTTP requests, executes it and
	// returns the response.
	Do(req *http.Request, opt Options) (*http.Response, error)
}
