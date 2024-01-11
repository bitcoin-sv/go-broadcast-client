package arcutils

import (
	"encoding/json"
	"io"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

func DecodeResponseBody(body io.ReadCloser, resultOutput any) error {
	err := json.NewDecoder(body).Decode(resultOutput)
	if err != nil {
		return broadcast.ErrUnableToDecodeResponse
	}

	return nil
}
