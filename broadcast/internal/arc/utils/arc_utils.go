package arcutils

import (
	"encoding/json"
	"errors"
	"fmt"
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

func WithCause(err error, cause error) error {
	return errors.Join(err, fmt.Errorf("\tcaused by: %w\t", cause))
}
