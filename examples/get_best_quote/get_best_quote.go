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

	taalCfg := broadcast_client.ArcClientConfig{
		Token:  token,
		APIUrl: apiURL,
	}

	gorillaCfg := broadcast_client.ArcClientConfig{
		Token:  "",
		APIUrl: "https://arc.gorillapool.io",
	}

	client := broadcast_client.Builder().
		WithArc(taalCfg).
		WithArc(gorillaCfg).
		Build()

	bestQuote, err := client.GetBestQuote(context.Background())
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	broadcast.PrettyPrint("Best Quote", bestQuote)
}
