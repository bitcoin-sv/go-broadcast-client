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
	hex := "469dd0f63fe4b3fe972dc72d28057e931abd69d8dfc85bf6e623efa41d50ef73"

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
