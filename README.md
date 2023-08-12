# go-broadcast-client

> Interact with Bitcoin SV Overlay Nodes exposing the interface [interface.go](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/interface.go)

## Features

- Arc API Support [details](https://github.com/bitcoin-sv/arc):
  - [x] [Query Transaction Status](https://bitcoin-sv.github.io/arc/api.html#/Arc/GET%20transaction%20status)
  - [x] [Submit Transaction](https://bitcoin-sv.github.io/arc/api.html#/Arc/POST%20transaction)
  - [x] [Submit Batch Transactions](https://bitcoin-sv.github.io/arc/api.html#/Arc/POST%20transactions)
  - [ ] Quote Services -> WORK IN PROGRESS

## What is library doing?

It is a wrapper around the [Bitcoin SV Overlay API](https://bitcoin-sv.github.io/arc/api.html) that allows you to interact with the API in a more convenient way by providing a set of
custom features to work with multiple nodes and retry logic.

## Custom features

- [x] Possibility to work with multiple nodes [builder pattern](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/broadcast-client/client_builder.go)
  
- [x] Define strategy how to work with multiple nodes [strategy](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/internal/composite/strategy.go)
  
- [x] Gives possibility to handle different client exposing the same interface as Arc [WithArc](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/broadcast-client/client_builder.go)
  
- [x] Possibility to set url and access token for each node independently [Config](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/broadcast-client/arc_config.go)
  
- [x] Possibility to use custom http client [WithHTTPClient](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/broadcast-client/client_builder.go#L19)

## How to use it?

### Create client

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

### Use the method exposed by the interface

```go
    // ...
    hex := "9c5f5244ee45e8c3213521c1d1d5df265d6c74fb108961a876917073d65fef14"

    result, err := client.QueryTransaction(context.Background(), hex)
    // ...
```

Examples of usage can be found in the [examples](https://github.com/bitcoin-sv/go-broadcast-client/tree/main/examples)

## ClientBuilder Methods

Client allows you to create your own client with setting it with multiple nodes and custom http client.

### WithArc Method

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

We can also call multiple times the method `WithArc` to set multiple nodes.

```go
// Set the access token and url for the node
 client := broadcast_client.Builder().
  WithArc(cfg1).
  WithArc(cfg2).
  WithArc(cfg3).
  Build()
```

What is the call order if we have multiple nodes configured?

We use the **strategy** to define the order of the calls. The default strategy is `OneByOne` in RoundRobin algorith that will call the nodes in the order they were set.
  **Only if all of them fail, we will return the error to the user.**

### WithHTTPClient Method

```go
// (...)
 clent := broadcast_client.Builder().
  WithArc(cfg).
  WithHTTPClient(&http.Client{}).
  Build()
```

We can use the method `WithHTTPClient` to set the custom http client.
It needs to implement the interface `HTTPInterface` that is defined in the [httpclient.go](/broadcast/internal/httpclient/http_client.go#L35)

```go
type HTTPInterface interface {
 DoRequest(ctx context.Context, pld HTTPRequest) (*http.Response, error)
}
```


## ARC Client Methods

### QueryTx Method

When you created your own client, you can use the method `QueryTx` to query the status of a transaction.

```go
    // ...
    hex := "9c5f5244ee45e8c3213521c1d1d5df265d6c74fb108961a876917073d65fef14"

    result, err := client.QueryTransaction(context.Background(), hex)
    // ...
```

### SubmitTx Method

Having your client created, you can use the method `SubmitTx` to submit a single transaction to the node.

```go
    // ...
    tx := broadcast.Transaction{
      RawTx: "xyz",
    }

    result, err := client.SubmitTransaction(context.Background(), tx)
    // ...
```

You need to pass the [transaction](#transaction) as a parameter to the method `SubmitTransaction`.

Setting tx.MerkleProof to true will add the header `X-MerkleProof` to the request.
MerkleProof while broadcasting will handle the merkle proof capability of the node.

Setting tx.CallBackURL and tx.CallBackToken will add the headers `X-CallbackUrl` and `X-CallbackToken` to the request.
It will allow you to get the callback from the node when the transaction is mined, and receive the transaction details and status.

Setting tx.WaitForStatus will add the header `X-WaitForStatus` to the request.
It will allow you to wait for the transaction to be mined and return the result only when the transaction reaches the status you set.

### SubmitBatchTx Method

Having your client created, you can use the method `SubmitBatchTx` to submit a batch of transactions to the node.

```go
    // ...
     txs := []*broadcast.Transaction{
      {RawTx: "xyz1"},
      {RawTx: "xyz2"},
      {RawTx: "xyz3"},
      {RawTx: "xyz4"},
      {RawTx: "xyz5"},
    }

    result, err := client.SubmitBatchTransaction(context.Background(), txs)
    // ...
```

The method works the same as the `SubmitTx` method, but it is sending a batch of transactions instead of a single one. It is also receiving a batch of responses for each of the transactions.

## Models and constants

### QueryTx

#### QueryTxResponse

```go
type QueryTxResponse struct {
 BlockHash string `json:"blockHash,omitempty"`
 BlockHeight int64 `json:"blockHeight,omitempty"`
 Timestamp string `json:"timestamp,omitempty"`
 TxID string `json:"txid,omitempty"`
 TxStatus TxStatus `json:"txStatus,omitempty"`
}
```

### SubmitTx

#### SubmitTxResponse

```go
type SubmitTxResponse struct {
 BlockHash string `json:"blockHash,omitempty"`
 BlockHeight int64 `json:"blockHeight,omitempty"`
 ExtraInfo string `json:"extraInfo,omitempty"`
 Status int `json:"status,omitempty"`
 Title string `json:"title,omitempty"`
 TxStatus TxStatus `json:"txStatus,omitempty"`
}
```

#### Transaction

  ```go
type Transaction struct {
  CallBackEncryption string `json:"callBackEncryption,omitempty"`
  CallBackToken string `json:"callBackToken,omitempty"`
  CallBackURL string `json:"callBackUrl,omitempty"`
  DsCheck bool `json:"dsCheck,omitempty"`
  MerkleFormat string `json:"merkleFormat,omitempty"`
  MerkleProof bool `json:"merkleProof,omitempty"`
  RawTx string `json:"rawtx"`
  WaitForStatus TxStatus `json:"waitForStatus,omitempty"`
}
  ```

### Transaction Statuses returned by the library

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
