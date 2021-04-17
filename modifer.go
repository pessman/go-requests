package requests

import "net/http"

const (
	applicationJson = "application/json"
)

type Modifier func(req *http.Request)

func DefaultModifer(req *http.Request) {
	req.Header.Set("Content-Type", applicationJson)
	req.Header.Set("Accept", applicationJson)
}
