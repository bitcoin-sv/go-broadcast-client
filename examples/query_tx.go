package main

import (
	"fmt"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
)

func main() {
	token := ""
	apiURL := ""

	cfg := broadcast_client.ArcClientConfig{
		Token:  token,
		APIUrl: apiURL,
	}

	client := broadcast_client.Builder().
		WithArc(cfg).
		Build()

	fmt.Print(client)
}
