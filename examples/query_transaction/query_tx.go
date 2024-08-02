package main

import (
	"context"
	"log"

	"github.com/rs/zerolog"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	broadcast_client "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
)

// This example shows how to query a transaction by its hex.
func main() {
	logger := zerolog.Nop()
	token := "mainnet_06770f425eb00298839a24a49cbdc02c"
	apiURL := "https://arc.taal.com"
	deploymentID := "broadcast-client-example"
	txID := "c67d504ea27eb225655bdb1f29c3a1a8f2056d675f3a7b17e91638efa30e1665"

	cfg := broadcast_client.ArcClientConfig{
		Token:        token,
		APIUrl:       apiURL,
		DeploymentID: deploymentID,
	}

	client := broadcast_client.Builder().
		WithArc(cfg, &logger).
		Build()

	result, err := client.QueryTransaction(context.Background(), txID)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	broadcast.PrettyPrint("Result: ", result)
}
