package errors

import "errors"

var ErrClientUndefined = errors.New("client is undefined")

var ErrAllBroadcastersFailed = errors.New("all broadcasters failed")

var ErrURLEmpty = errors.New("url is empty")

var ErrBroadcasterFailed = errors.New("broadcaster failed")

var ErrUnableToParseResponse = errors.New("unable to parse response")
