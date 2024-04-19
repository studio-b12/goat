package requester

// Options wraps request specific options.
type Options struct {
	CookieJar       any
	StoreCookies    bool
	SendCookies     bool
	ResponseType    string
	FollowRedirects bool
}

// OptionsFromMap takes a map and builds an
// Options instance from matching key-value
// pairs.
func OptionsFromMap(m map[string]any) Options {
	opt := Options{
		CookieJar:       "default",
		StoreCookies:    true,
		SendCookies:     true,
		ResponseType:    "",
		FollowRedirects: true,
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
	if v, ok := m["responsetype"].(string); ok {
		opt.ResponseType = v
	}
	if v, ok := m["followredirects"].(bool); ok {
		opt.FollowRedirects = v
	}

	return opt
}
