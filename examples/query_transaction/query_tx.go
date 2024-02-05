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
	token := ""
	apiURL := "https://arc.gorillapool.io"
	deploymentID := "broadcast-client-example"
	txID := "469dd0f63fe4b3fe972dc72d28057e931abd69d8dfc85bf6e623efa41d50ef73"

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
