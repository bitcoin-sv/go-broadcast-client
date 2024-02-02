package main

import (
	"context"
	"log"

	"github.com/rs/zerolog"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	broadcast_client "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
)

// This example shows how to get a fee quote.
func main() {
	logger := zerolog.Nop()

	gorillaCfg := broadcast_client.ArcClientConfig{
		Token:  "",
		APIUrl: "https://arc.gorillapool.io",
	}

	client := broadcast_client.Builder().
		WithArc(gorillaCfg, &logger, broadcast_client.WithXDeploymentID("broadcast-client-example")).
		Build()

	feeQuotes, err := client.GetFeeQuote(context.Background())
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	broadcast.PrettyPrint("Fee Quotes", feeQuotes)
}
