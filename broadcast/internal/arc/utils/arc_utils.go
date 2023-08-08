package arc_utils

import (
	"encoding/json"
	"io"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

func DecodeArcError(body io.ReadCloser) error {
	resultError := broadcast.ArcError{}
	err := json.NewDecoder(body).Decode(&resultError)

	if err != nil {
		return broadcast.ErrUnableToDecodeResponse
	}

	return resultError
}

func DecodeResponseBody(body io.ReadCloser, resultOutput any) error {
	err := json.NewDecoder(body).Decode(resultOutput)
	if err != nil {
		return broadcast.ErrUnableToDecodeResponse
	}

	return nil
}
