package main

import (
	"context"
	"log"

	"github.com/rs/zerolog"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	broadcast_client "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
)

// This example shows how to submit a transaction.
func main() {
	logger := zerolog.Nop()
	token := "mainnet_06770f425eb00298839a24a49cbdc02c"
	deploymentID := "broadcast-client-example"
	apiURL := "https://arc.taal.com"
	tx := broadcast.Transaction{Hex: "0100000001f58d6ead84eadf358f85db9347e687ec5ed15559330e6b0929e967054ff5dbc4000000006a47304402207831d6ee02ebda30ff306946a1e63654d8234a99e217eb5ddb122fec10b9ea0a02204d18c24990df7fd4efe6805869389172921b473e34a31e35b336c402085f9b964121038989293fb3a740c239977c42edc679831b7f6cb1005545b608034a85f114baf3ffffffff020100000000000000084c0105954c01198807270000000000001976a914e9e47358fc4244d146161527c1ae870bcc179b3488ac00000000"}

	cfg := broadcast_client.ArcClientConfig{
		Token:        token,
		APIUrl:       apiURL,
		DeploymentID: deploymentID,
	}

	client := broadcast_client.Builder().
		WithArc(cfg, &logger).
		Build()

	result, err := client.SubmitTransaction(context.Background(), &tx, broadcast.WithRawFormat())
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	broadcast.PrettyPrint("Result: ", result)
}
