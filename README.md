# WORK IN PROGRESS

# go-broadcast-client
> Interact with Bitcoin SV Overlay Nodes exposing the interface [interface.go](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/interface.go)


# Features
- Arc API Support [details](https://github.com/bitcoin-sv/arc):
  - [x] [Query Transaction Status](https://bitcoin-sv.github.io/arc/api.html#/Arc/GET%20transaction%20status)
  - [x] [Submit Transaction](https://bitcoin-sv.github.io/arc/api.html#/Arc/POST%20transaction)
  - [ ] Quote Services -> WORK IN PROGRESS
  - [x] [Submit Batch Transactions](https://bitcoin-sv.github.io/arc/api.html#/Arc/POST%20transactions)

# What is library doing?
It is a wrapper around the [Bitcoin SV Overlay API](https://bitcoin-sv.github.io/arc/api.html) that allows you to interact with the API in a more convenient way by providing a set of
custom features to work with multiple nodes and retry logic.

# Custom features
  - [x] Possibility to work with multiple nodes [builder pattern](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/broadcast-client/client_builder.go)
  - [x] Define strategy how to work with multiple nodes [strategy](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/internal/composite/strategy.go)
  - [x] Gives possibility to handle different client exposing the same interface as Arc [WithArc](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/broadcast-client/client_builder.go)
  - [x] Possibility to set url and access token for each node independently [Config](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/broadcast-client/arc_config.go)
  - [x] Possibility to use custom http client [WithHTTPClient](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/broadcast-client/client_builder.go#L19)

# How to use it?

## Create client
```go
    // Set the access token and url for the node
	token := ""
	apiURL := "https://tapi.taal.com/arc"

	cfg := broadcast_client.ArcClientConfig{
		Token:  token,
		APIUrl: apiURL,
	}

	client := broadcast_client.Builder().
		WithArc(cfg).
		Build()
```

## Use the method exposed by the interface
```go
    // ...
    hex := "9c5f5244ee45e8c3213521c1d1d5df265d6c74fb108961a876917073d65fef14"

    result, err := client.QueryTransaction(context.Background(), hex)
    // ...
```

Examples of usage can be found in the [examples](https://github.com/bitcoin-sv/go-broadcast-client/tree/main/examples)

# Transaction Statuses returned by the library

| Code | Status                 | Description                                                                                                                                                                                                    |
|-----|------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 0   | `UNKNOWN`              | The transaction has been sent to metamorph, but no processing has taken place. This should never be the case, unless something goes wrong.                                                                     |
| 1   | `QUEUED`               | The transaction has been queued for processing.                                                                                                                                                                |
| 2   | `RECEIVED`             | The transaction has been properly received by the metamorph processor.                                                                                                                                         |
| 3   | `STORED`               | The transaction has been stored in the metamorph store. This should ensure the transaction will be processed and retried if not picked up immediately by a mining node.                                        |
| 4   | `ANNOUNCED_TO_NETWORK` | The transaction has been announced (INV message) to the Bitcoin network.                                                                                                                                       |
| 5   | `REQUESTED_BY_NETWORK` | The transaction has been requested from metamorph by a Bitcoin node.                                                                                                                                           |
| 6   | `SENT_TO_NETWORK`      | The transaction has been sent to at least 1 Bitcoin node.                                                                                                                                                      |
| 7   | `ACCEPTED_BY_NETWORK`  | The transaction has been accepted by a connected Bitcoin node on the ZMQ interface. If metamorph is not connected to ZQM, this status will never by set.                                                       |
| 8   | `SEEN_ON_NETWORK`      | The transaction has been seen on the Bitcoin network and propagated to other nodes. This status is set when metamorph receives an INV message for the transaction from another node than it was sent to.       |
| 9   | `MINED`                | The transaction has been mined into a block by a mining node.                                                                                                                                                  |
| 108 | `CONFIRMED`            | The transaction is marked as confirmed when it is in a block with 100 blocks built on top of that block.                                                                                                       |
| 109 | `REJECTED`             | The transaction has been rejected by the Bitcoin network.   

*Source* [Arc API](https://github.com/bitcoin-sv/arc/blob/main/README.md)
