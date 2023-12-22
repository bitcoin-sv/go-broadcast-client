package main

import (
	"context"
	"github.com/rs/zerolog"
	"log"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	broadcast_client "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
)

// This example shows how to get a fee quote.
func main() {
	logger := zerolog.Nop()
	taalCfg := broadcast_client.ArcClientConfig{
		Token:  "",
		APIUrl: "https://tapi.taal.com/arc",
	}

	gorillaCfg := broadcast_client.ArcClientConfig{
		Token:  "",
		APIUrl: "https://arc.gorillapool.io",
	}

	client := broadcast_client.Builder().
		WithArc(taalCfg, &logger).
		WithArc(gorillaCfg, &logger).
		Build()

	feeQuotes, err := client.GetFeeQuote(context.Background())
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	broadcast.PrettyPrint("Fee Quotes", feeQuotes)
}
