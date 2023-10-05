package main

import (
	"context"
	"log"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	broadcast_client "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
)

// This example shows how to submit a transaction.
func main() {
	token := ""
	apiURL := "https://tapi.taal.com/arc"
	tx := broadcast.Transaction{Hex: "0100000001d6d1607b208b30c0a3fe21d563569c4d2a0f913604b4c5054fe267da6be324ab220000006b4830450221009a965dcd5d42983090a63cfd761038ff8adcea621c46a68a205f326292a95383022061b8d858f366c69f3ebd30a60ccafe36faca4e242ac3d2edd3bf63b669bcf23b4121034e871e147aa4a3e2f1665eaf76cf9264d089b6a91702af92bd6ce33bac84a765ffffffff0123020000000000001976a914d8819a7197d3e221e15f4348203fdecfd29fa2b888ac00000000"}

	cfg := broadcast_client.ArcClientConfig{
		Token:  token,
		APIUrl: apiURL,
	}

	client := broadcast_client.Builder().
		WithArc(cfg).
		Build()

	result, err := client.SubmitTransaction(context.Background(), &tx, broadcast.WithRawFormat())
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	log.Printf("miner: %s", result.Miner)
	log.Printf("hash: %s", result.BlockHash)
	log.Printf("status: %s", result.TxStatus)
}
