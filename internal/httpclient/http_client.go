package httpclient

import (
	"io"
	"net/http"

	errors "github.com/bitcoin-sv/go-broadcast-client"
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
)

type HTTPClient struct {
	Client *http.Client
}

type HTTPInterface interface {
	DoRequest(method HttpMethod, url, token string, body io.Reader) (*http.Response, error)
}

func NewHttpClient() HTTPInterface {
	return &HTTPClient{
		Client: &http.Client{},
	}
}

func (hc *HTTPClient) DoRequest(method HttpMethod, url, token string, body io.Reader) (*http.Response, error) {
	if url == "" {
		return nil, errors.ErrURLEmpty
	}

	req, err := http.NewRequest(string(method), url, body)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Add("Authorization", token)
	}

	resp, err := hc.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
