package client

import "net/http"

// Client lets add a comment to test
type Client interface {
	Do(*http.Request) (*http.Response, error)
}
