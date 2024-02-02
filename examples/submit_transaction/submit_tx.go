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
	token := ""
	apiURL := "https://arc.gorillapool.io"
	tx := broadcast.Transaction{Hex: "01000000026d9a56036235bb5b5e39b04b6f188c74bc45189a9122f235cad5f0c4b668817d000000006b483045022100ae4c9b06376c42bf82f7910de20bb025d8b43daf33eb2db1966f0f1fd361d499022063594799502920ceceb17e3301e44066431a5ae1e221ce1bd89b446e602adf62412102dd2063cc1d4fbc5a770b156f8cd9f5c80cc586df4c7d148444d6bb66c81a10daffffffff6e574c52ebc1c724a13d0afd525adbcfe0b134d9666ab9004240d9d072a7906b000000006a4730440220191c938793953f931c9297a27a651b07e5cb60432fd5750a3f53488ed23fae8702204c503822986ef959af741609dfe98a51004e36596ca52a865c5e029fdbcfb3d641210253085022df5ebbdc71f9e8f555443660cda4a3d36d732685be317d7e11b9eda4ffffffff0214000000000000001976a914553236189ff9fed552837b952e404e09f78c03fa88ac01000000000000001976a914d8f00ced3f960ffa6f0516b4e352e0a9feccc3af88ac00000000"}

	cfg := broadcast_client.ArcClientConfig{
		Token:  token,
		APIUrl: apiURL,
	}

	client := broadcast_client.Builder().
		WithArc(cfg, &logger, broadcast_client.WithXDeploymentID("broadcast-client-example")).
		Build()

	result, err := client.SubmitTransaction(context.Background(), &tx, broadcast.WithRawFormat())
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	broadcast.PrettyPrint("Result: ", result)
}
