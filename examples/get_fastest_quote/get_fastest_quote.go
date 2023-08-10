package main

import (
	"context"
	"log"
	"time"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	broadcast_client "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
)

func main() {
	taalCfg := broadcast_client.ArcClientConfig{
		Token:  "",
		APIUrl: "https://tapi.taal.com/arc",
	}

	gorillaCfg := broadcast_client.ArcClientConfig{
		Token:  "",
		APIUrl: "https://arc.gorillapool.io",
	}

	client := broadcast_client.Builder().
		WithArc(gorillaCfg).
		WithArc(taalCfg).
		Build()

	bestQuote, err := client.GetFastestQuote(context.Background(), time.Second)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	broadcast.PrettyPrint("Best Quote", bestQuote)
}
