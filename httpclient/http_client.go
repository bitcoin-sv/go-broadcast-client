package httpclient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
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

func genericErrorMessage(statusCode int, message string) error {
	return fmt.Errorf("server responded with no-success code. details: { statusCode: %d, body: %s }", statusCode, message)
}

func handleErrorResponse(errorResponse *http.Response, errorParser ErrorParserFunction) error {
	body, err := io.ReadAll(errorResponse.Body)
	if err != nil {
		return fmt.Errorf("error during reading body: %s", err.Error())
	}

	if errorParser != nil {
		err = errorParser(errorResponse.StatusCode, body)
		if err != nil {
			return err
		}
	}

	return genericErrorMessage(errorResponse.StatusCode, string(body))
}

func (hc *HTTPClient) DoRequest(ctx context.Context, pld HTTPRequest) (*http.Response, error) {
	var bodyReader io.Reader

	if pld.URL == "" {
		return nil, errors.New("url is empty")
	}

	if pld.Data != nil && (pld.Method == http.MethodPost || pld.Method == http.MethodPut) {
		bodyReader = bytes.NewBuffer(pld.Data)
	}

	req, err := http.NewRequestWithContext(ctx, string(pld.Method), pld.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if pld.Token != "" {
		req.Header.Add("Authorization", pld.Token)
	}

	if pld.Headers != nil {
		for key, value := range pld.Headers {
			req.Header.Add(key, value)
		}
	}

	return hc.Client.Do(req)
}

func hasSuccessCode(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

type RequestFuncion func(ctx context.Context, pld HTTPRequest) (*http.Response, error)
type ParserFunction[T any] func(resp *http.Response) (T, error)
type ErrorParserFunction func(statusCode int, body []byte) error

func RequestModel[T any](
	ctx context.Context,
	req RequestFuncion,
	pld HTTPRequest,
	parser ParserFunction[T],
	errorParser ErrorParserFunction,
) (model T, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("requestModel panic: %v", r)
		}
	}()
	resp, err := req(ctx, pld)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}

	if !hasSuccessCode(resp) {
		err = handleErrorResponse(resp, errorParser)
		return
	}

	model, err = parser(resp)
	if err != nil {
		return
	}

	return
}
