package arc

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

func supportedArcError(statusCode int) bool {
	switch statusCode {
	case 400, 422, 465, 466:
		return true
	default:
		return false
	}
}

func parseArcError(statusCode int, body []byte) error {
	if !supportedArcError(statusCode) {
		return nil
	}

	resultError := broadcast.ArcError{}
	reader := bytes.NewReader(body)
	err := json.NewDecoder(reader).Decode(&resultError)

	if err != nil {
		return errors.New("unable to decode arc error")
	}

	if resultError.Title != "" {
		return &resultError
	}

	return nil
}
