package main

import (
	"fmt"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

func main() {
	token := ""
	apiURL := ""

	cfg := broadcast.ArcClientConfig{
		Token:  token,
		APIUrl: apiURL,
	}

	client := broadcast.NewClientBuilder().
		WithArc(cfg).
		Build()

	fmt.Print(client)
}
