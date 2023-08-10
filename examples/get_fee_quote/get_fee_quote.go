package main

import (
	"context"
	"log"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	broadcast_client "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
)

func main() {
	token := ""
	apiURL := "https://tapi.taal.com/arc"

	cfg := broadcast_client.ArcClientConfig{
		Token:  token,
		APIUrl: apiURL,
	}

	client := broadcast_client.Builder().
		WithArc(cfg).
		Build()

	feeQuote, err := client.GetFeeQuote(context.Background())
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	broadcast.PrettyPrint("Fee Quote", feeQuote)
}
