package config

import "github.com/bitcoin-sv/go-broadcast-client/broadcast"

type Config interface{}

var DefaultStrategy broadcast.Strategy = *broadcast.OneByOne
