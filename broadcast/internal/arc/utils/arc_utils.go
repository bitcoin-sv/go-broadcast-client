package arcutils

import (
	"encoding/json"
	"io"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

func DecodeResponseBody[T any](body io.ReadCloser, resultOutput *T) error {
	err := json.NewDecoder(body).Decode(resultOutput)
	if err != nil {
		return broadcast.ErrUnableToDecodeResponse
	}

	return nil
}
