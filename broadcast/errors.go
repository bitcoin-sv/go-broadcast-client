package broadcast

import (
	"errors"
	"fmt"
	"strings"
)

var ErrClientUndefined = errors.New("client is undefined")

var ErrAllBroadcastersFailed = errors.New("all broadcasters failed")

var ErrURLEmpty = errors.New("url is empty")

var ErrBroadcasterFailed = errors.New("broadcaster failed")

var ErrUnableToDecodeResponse = errors.New("unable to decode response")

var ErrMissingStatus = errors.New("missing tx status")

var ErrStrategyUnkown = errors.New("unknown strategy")

type ArcError struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Status    int    `json:"status"`
	Detail    string `json:"detail"`
	Instance  string `json:"instance,omitempty"`
	Txid      string `json:"txid,omitempty"`
	ExtraInfo string `json:"extraInfo,omitempty"`
}

func (err ArcError) Error() string {
	sb := strings.Builder{}

	sb.WriteString("arc error: {")
	sb.WriteString(fmt.Sprintf("type: %s, title: %s, status: %d, detail: %s",
		err.Type, err.Title, err.Status, err.Detail))

	if err.Instance != "" {
		sb.Write([]byte(fmt.Sprintf(", instance: %s", err.Instance)))
	}

	if err.Txid != "" {
		sb.Write([]byte(fmt.Sprintf(", txid: %s", err.Txid)))
	}

	if err.ExtraInfo != "" {
		sb.Write([]byte(fmt.Sprintf(", extraInfo: %s", err.ExtraInfo)))
	}

	sb.WriteString("}")
	return sb.String()
}
