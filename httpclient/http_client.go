package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

func errorDuringReadingBody(err error) string {
	return fmt.Sprintf("error during reading body: %s", err.Error())
}

func decodeErrorResponse(errorResponse *http.Response) error {
	body, err := io.ReadAll(errorResponse.Body)

	var message string
	if err != nil {
		message = fmt.Sprintf("error during reading body: %s", err.Error())
	} else {
		message = string(body)
	}
	return genericErrorMessage(errorResponse.StatusCode, message)
}

func decodeArcError(errorResponse *http.Response) error {
	body, err := io.ReadAll(errorResponse.Body)
	if err != nil {
		return errors.New(errorDuringReadingBody(err))
	}

	resultError := broadcast.ArcError{}
	reader := bytes.NewReader(body)
	err = json.NewDecoder(reader).Decode(&resultError)

	if err != nil {
		return errors.New("unable to decode arc error")
	}

	if resultError.Title != "" {
		return resultError
	}

	return genericErrorMessage(errorResponse.StatusCode, string(body))
}

func handleErrorResponse(errorResponse *http.Response) error {
	switch errorResponse.StatusCode {
	case 400, 422, 465, 466:
		return decodeArcError(errorResponse)
	default:
		return decodeErrorResponse(errorResponse)
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

type requestFuncion func(ctx context.Context, pld HTTPRequest) (*http.Response, error)

func RequestModel[T any](ctx context.Context, req requestFuncion, pld HTTPRequest, parser func(resp *http.Response) (T, error)) (model T, err error) {
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
		err = handleErrorResponse(resp)
		return
	}

	model, err = parser(resp)
	if err != nil {
		return
	}

	return
}
