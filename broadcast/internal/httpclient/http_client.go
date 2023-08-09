package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
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

type HttpClientError struct {
	Response *http.Response
}

func (err HttpClientError) Error() string {
	body, _ := io.ReadAll(err.Response.Body)
	return fmt.Sprintf("server responded with no-success code. details: { statusCode: %d, body: %s }", err.Response.StatusCode, body)
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
		return nil, broadcast.ErrURLEmpty
	}

	req, err := http.NewRequestWithContext(ctx, string(pld.Method), pld.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	if pld.Data != nil && (pld.Method == POST || pld.Method == PUT) {
		bodyReader = bytes.NewBuffer(pld.Data)
		req.Body = io.NopCloser(bodyReader)
		req.Header.Set("Content-Type", "application/json")
	}

	if pld.Token != "" {
		req.Header.Add("Authorization", pld.Token)
	}

	resp, err := hc.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if hasSuccessCode(resp) {
		return resp, nil
	}

	return nil, HttpClientError{resp}
}

func hasSuccessCode(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
