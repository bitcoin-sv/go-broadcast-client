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
	apiURL := "https://api.taal.com/arc"

	//hex := "0100000001c73ec4cebb3f4fbb4eba6fd8e481d0e50ec004e74abc160ab595b0797222ed04010000006b483045022100edf7220bd6ba442180593d618bc22c863756cd759912412e8d4affff064120ba0220206c97625bcedb19f5edb6f063da5ed1d07f8aa4ad7c5885e304760a1c0357674121028725265e99175bd3c5f356ab93415b6d8da75d6c23da5f2ac49990db8b8e6bb5ffffffff0200000000000000000a006a07616b756b7520365d000000000000001976a9142f1609de754218a8489b6d021ebc9d0a7d35802088ac00000000"
	hex := "0100000002f52bcbe41d5339b08c81bf8f50312cb164867a5d6ffb54465570332b02c05330010000006a4730440220222c63eb9b7b84f8812cc09dae554da4aecb868b26c33bcd06b8ac6d1cf7868702202a7a41326d3d10dfb059c25700722863db9be730c4df077d089aa37f3db3a55d412102c7a5a1b2e739ad5ef9e4d1626151208bf4aefe27cff8cb0b6ec8ef4f1cf702a8ffffffff0200000000000000000a006a07616b756b7520355e000000000000001976a914c45076ef1b386d75f1036b8225fc2715ca8f267d88ac00000000"
	tx := broadcast.Transaction{Hex: hex}

	cfg := broadcast_client.ArcClientConfig{
		Token:  token,
		APIUrl: apiURL,
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
