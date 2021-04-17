package requests

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pessman/go-requests/client"
	"github.com/pessman/go-requests/response"
)

type Request interface {
	Response() (*http.Response, error)
	BytesResponse() (*response.BytesResponse, error)
	// ReaderResponse() (io.Reader, error)
}

type request struct {
	URL       string
	Method    string
	Body      io.Reader
	Client    client.Client
	Modifiers []Modifier
}

func (r *request) Response() (*http.Response, error) {
	if r.Client == nil {
		r.Client = http.DefaultClient
	}

	req, err := r.newHTTPRequest()
	if err != nil {
		return nil, err
	}

	return r.Client.Do(req)
}

func (r *request) BytesResponse() (*response.BytesResponse, error) {
	res, err := r.Response()
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return response.NewBytesResponse(b, res.StatusCode), nil
}

func (r *request) newHTTPRequest() (*http.Request, error) {
	req, err := http.NewRequest(r.Method, r.URL, r.Body)
	if err != nil {
		return nil, err
	}

	for _, m := range r.Modifiers {
		m(req)
	}

	return req, nil
}

func NewRequest(url, method string, body io.Reader, client client.Client) Request {
	return &request{
		URL:    url,
		Method: method,
		Body:   body,
		Client: client,
	}
}
