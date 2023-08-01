package main

import (
	"context"
	"log"

	broadcast_client "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
)

func main() {
	token := ""
	apiURL := "https://tapi.taal.com/arc"
	hex := "9c5f5244ee45e8c3213521c1d1d5df265d6c74fb108961a876917073d65fef14"

	cfg := broadcast_client.ArcClientConfig{
		Token:  token,
		APIUrl: apiURL,
	}

	client := broadcast_client.Builder().
		WithArc(cfg, nil).
		Build()

	result, err := client.QueryTransaction(context.Background(), hex)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	log.Printf("hash: %s", result.BlockHash)
	log.Printf("status: %s", result.TxStatus)
	log.Printf("block height: %d", result.BlockHeight)
}
