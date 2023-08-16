package arcutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/httpclient"
)

var ErrUnableToDecodeArcError = errors.New("unable to decode arc error")

func HandleHttpError(httpClientError error) error {
	noSuccessResponseErr, ok := httpClientError.(httpclient.HttpClientError)

	if ok { // client respond with code different than 2xx
		var err error

		switch noSuccessResponseErr.Response.StatusCode {
		case 400:
			err = decodeArcError(noSuccessResponseErr)
		case 422: // 	Unprocessable entity - with IETF RFC 7807 Error object
			err = decodeArcError(noSuccessResponseErr)
		case 465: // 	Fee too low
			err = decodeArcError(noSuccessResponseErr)
		case 466: // 	Conflicting transaction found
			err = decodeArcError(noSuccessResponseErr)

		default:
			err = noSuccessResponseErr
		}

		return err
	}

	return httpClientError // http client internal error
}

func DecodeResponseBody(body io.ReadCloser, resultOutput any) error {
	err := json.NewDecoder(body).Decode(resultOutput)
	if err != nil {
		return broadcast.ErrUnableToDecodeResponse
	}

	return nil
}

func decodeArcError(httpErr httpclient.HttpClientError) error {
	response := httpErr.Response
	// duplicate stream
	var buffer bytes.Buffer
	bodyReader := io.TeeReader(response.Body, &buffer)

	resultError := broadcast.ArcError{}
	err := json.NewDecoder(bodyReader).Decode(&resultError)

	if err != nil {
		return ErrUnableToDecodeArcError
	}

	if resultError.Title != "" {
		return resultError
	}

	// miner returns an error with an invalid schema
	return httpErr
}
