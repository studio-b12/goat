package requester

import "net/http"

// Requester defines a service to perform HTTP
// requests.
type Requester interface {
	// Do takes a HTTP requests, executes it and
	// returns the response.
	Do(req *http.Request) (*http.Response, error)
}
