// Package httpclient provides the custom implementation of the http client.
package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

// HttpMethod is a custom type for the http method.
type HttpMethod string

const (
	// GET is the http method for GET requests.
	GET HttpMethod = "GET"
	// POST is the http method for POST requests.
	POST HttpMethod = "POST"
	// PUT is the http method for PUT requests.
	PUT HttpMethod = "PUT"
	// DELETE is the http method for DELETE requests.
	DELETE HttpMethod = "DELETE"
)

// HTTPClient is the custom implementation of the http client.
type HTTPClient struct {
	// Client is the http client.
	Client *http.Client
}

// HTTPInterface is the interface for the http client.
type HTTPInterface interface {
	// DoRequest performs the http request.
	DoRequest(ctx context.Context, pld HTTPRequest) (*http.Response, error)
}

// HTTPRequest is the custom implementation of payload of the http request.
type HTTPRequest struct {
	// Method is the http method.
	Method HttpMethod
	// URL is the url of the http request.
	URL string
	// Token is the token of the http request.
	Token string
	// Data is the data of the http request.
	Data []byte
	// Headers is the map of the headers of the http request.
	Headers map[string]string
}

// HttpClientError is the custom implementation of the http client error.
type HttpClientError struct {
	// Response is the http response.
	Response *http.Response
}

// Error returns the error message.
func (err HttpClientError) Error() string {
	body, _ := io.ReadAll(err.Response.Body)
	return fmt.Sprintf("server responded with no-success code. details: { statusCode: %d, body: %s }", err.Response.StatusCode, body)
}

// AddHeader adds the header to the http request.
func (pld *HTTPRequest) AddHeader(key, value string) {
	if pld.Headers == nil {
		pld.Headers = make(map[string]string)
	}

	pld.Headers[key] = value
}

// NewPayload creates the new payload of the http request.
func NewPayload(method HttpMethod, url, token string, data []byte) HTTPRequest {
	return HTTPRequest{
		Method: method,
		URL:    url,
		Token:  token,
		Data:   data,
	}
}

// NewHttpClient creates the new http client.
func NewHttpClient() HTTPInterface {
	return &HTTPClient{
		Client: &http.Client{},
	}
}

// DoRequest performs the http request.
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
