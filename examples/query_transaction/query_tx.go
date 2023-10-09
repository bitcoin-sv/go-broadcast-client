package main

import (
	"context"
	"log"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	broadcast_client "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
)

// This example shows how to query a transaction by its hex.
func main() {
	token := ""
	apiURL := "https://tapi.taal.com/arc"
	hex := "680d975a403fd9ec90f613e87d17802c029d2d930df1c8373cdcdda2f536a1c0"

	cfg := broadcast_client.ArcClientConfig{
		Token:  token,
		APIUrl: apiURL,
	}

	client := broadcast_client.Builder().
		WithArc(cfg).
		Build()

	result, err := client.QueryTransaction(context.Background(), hex)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	broadcast.PrettyPrint("Result: ", result)
}
