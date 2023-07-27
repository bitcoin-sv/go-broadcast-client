package shared

import (
	"errors"
)

var ErrClientUndefined = errors.New("client is undefined")

var ErrAllBroadcastersFailed = errors.New("all broadcasters failed")

var ErrURLEmpty = errors.New("url is empty")

var ErrBroadcasterFailed = errors.New("broadcaster failed")

var ErrUnableToDecodeResponse = errors.New("unable to decode response")

var ErrMissingStatus = errors.New("missing tx status")
