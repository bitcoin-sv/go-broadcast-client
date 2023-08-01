package httpclient

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-api"
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
	DoRequest(ctx context.Context, pld HTTPRequest) (*http.Response, error)
}

type HTTPRequest struct {
	Method  HttpMethod
	URL     string
	Token   string
	Data    []byte
	Headers map[string]string
}

func (pld *HTTPRequest) AddHeader(key, value string) {
	if pld.Headers == nil {
		pld.Headers = make(map[string]string)
	}

	pld.Headers[key] = value
}

func NewPayload(method HttpMethod, url, token string, data []byte) HTTPRequest {
	return HTTPRequest{
		Method: method,
		URL:    url,
		Token:  token,
		Data:   data,
	}
}

func NewHttpClient() HTTPInterface {
	return &HTTPClient{
		Client: &http.Client{},
	}
}

func (hc *HTTPClient) DoRequest(ctx context.Context, pld HTTPRequest) (*http.Response, error) {
	var bodyReader io.Reader

	if pld.URL == "" {
		return nil, broadcast_api.ErrURLEmpty
	}

	req, err := http.NewRequestWithContext(ctx, string(pld.Method), pld.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	if pld.Data != nil && (pld.Method == POST || pld.Method == PUT) {
		bodyReader = bytes.NewBuffer(pld.Data)
		req.Header.Set("Content-Type", "application/json")
	}

	if pld.Token != "" {
		req.Header.Add("Authorization", pld.Token)
	}

	resp, err := hc.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
