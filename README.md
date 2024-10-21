⚠️ This repository is deprecated. ⚠️
We strongly recommend to use instead
- [go-sdk](https://github.com/bitcoin-sv/go-sdk) instead and the transaction broadcaster there, of coure feel free to create an [issue in go-sdk](https://github.com/bitcoin-sv/go-sdk/issues), if there is any lacking feature in go-sdk that you used from this library 
- plain requests to ARC or maybe look at the [ARC repository](https://github.com/bitcoin-sv/arc) for some help with creating a client

# go-broadcast-client

> Interact with Bitcoin SV ARC exposing the interface [interface.go](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/interface.go)

## Features

- Arc API Support [details](https://github.com/bitcoin-sv/arc):
  - [x] [Query Transaction Status](https://bitcoin-sv.github.io/arc/api.html#/Arc/GET%20transaction%20status)
  - [x] [Submit Transaction](https://bitcoin-sv.github.io/arc/api.html#/Arc/POST%20transaction)
  - [x] [Submit Batch Transactions](https://bitcoin-sv.github.io/arc/api.html#/Arc/POST%20transactions)
  - [x] [Quote Services](https://bitcoin-sv.github.io/arc/api.html#/Arc/GET%20policy)

## What is library doing?

It is a wrapper around the [Bitcoin SV ARC API](https://bitcoin-sv.github.io/arc/api.html) that allows you to interact with the API in a more convenient way by providing a set of
custom features to work with multiple nodes and retry logic.

## Custom features

- [x] Possibility to work with multiple nodes [builder pattern](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/broadcast-client/client_builder.go)
  
- [x] Define strategy how to work with multiple nodes [strategy](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/internal/composite/strategy.go)
  
- [x] Gives possibility to handle different client exposing the same interface as Arc [WithArc](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/broadcast-client/client_builder.go)
  
- [x] Possibility to set url and access token for each node independently [Config](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/broadcast-client/arc_config.go)
  
- [x] Possibility to use custom http client [WithHTTPClient](https://github.com/bitcoin-sv/go-broadcast-client/blob/main/broadcast/broadcast-client/client_builder.go#L19)

- [x] Mock Client for testing purposes [details](#MockClientBuilder)

## How to use it?

### Create client

```go
 // Set the access token and url for the node
 token := ""
 apiURL := "https://tapi.taal.com/arc"
 // also there is a possibility to set deploymentID which is optional and adds header XDeployment-ID to the request
 deploymentID := "broadcast-client-example"

 cfg := broadcast_client.ArcClientConfig{
    Token:  token,
    APIUrl: apiURL,
    DeploymentID: deploymentID,
 }

 client := broadcast_client.Builder().
  WithArc(cfg, &logger).
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

### Get Policy Quote Method

Having your client created, you can use the method `GetPolicyQuote` to get a slice of policy quotes from all the nodes configured in your client.

```go
// ...
policyQuotes, err := client.GetPolicyQuote(context.Background())
// ...
```

### Get Fee Quote Method

Having your client created, you can use the method `GetFeeQuote` to get a slice of fee quotes (which are subsets from PolicyQuotes) from all the nodes configured in your client.

```go
// ...
feeQuotes, err := client.GetFeeQuote(context.Background())
// ...
```

### SubmitTx Method

Having your client created, you can use the method `SubmitTx` to submit a single transaction to the node.

```go
    // ...
    tx := broadcast.Transaction{
      Hex: "xyz",
    }

    result, err := client.SubmitTransaction(context.Background(), tx)
    // ...
```

You need to pass the [transaction](#transaction) as a parameter to the method `SubmitTransaction`.

You may add options to this method:

##### WithCallback

```go
    result, err := client.SubmitTransaction(context.Background(), tx, broadcast.WithCallback(callBackUrl, callbackToken))
```
Setting `CallbackURL` and `CallBackToken` will add the headers `X-CallbackUrl` and `X-CallbackToken` to the request.
It will allow you to get the callback from the node when the transaction is mined, and receive the transaction details and status.

##### WithMerkleProof

```go
    result, err := client.SubmitTransaction(context.Background(), tx, broadcast.WithMerkleProof())
```
Setting `MerkleProof` to true will add the header `X-MerkleProof` to the request.
MerkleProof while broadcasting will handle the merkle proof capability of the node.

##### WithWaitForstatus

```go
    result, err := client.SubmitTransaction(context.Background(), tx, broadcast.WithWaitForstatus(broadcast.AnnouncedToNetwork))
```
Setting `WaitForStatus` will add the header `X-WaitForStatus` to the request.
It will allow you to return the result only when the transaction reaches the status you set.

##### WithEfFormat

```go
    result, err := client.SubmitTransaction(context.Background(), tx, broadcast.WithEfFormat())
```
Setting `EfFormat` will accept your transaction in EF format and submit to Arc. This is the **default** option.

##### WithBeefFormat

```go
    result, err := client.SubmitTransaction(context.Background(), tx, broadcast.WithBeefFormat())
```
Setting `BeefFormat` will accept your transaction in BEEF format and decode it for a proper format acceptable by Arc.

##### WithRawFormat (**DEPRECATED!**)

```go
    result, err := client.SubmitTransaction(context.Background(), tx, broadcast.WithRawFormat())
```
Setting `RawFormat` will accept your transaction in RawTx format and encode it for a proper format acceptable by Arc.
This option will become deprecated soon.

### SubmitBatchTx Method

Having your client created, you can use the method `SubmitBatchTx` to submit a batch of transactions to the node.

```go
    // ...
     txs := []*broadcast.Transaction{
      {Hex: "xyz1"},
      {Hex: "xyz2"},
      {Hex: "xyz3"},
      {Hex: "xyz4"},
      {Hex: "xyz5"},
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
    Miner       string         `json:"miner"`
    BlockHash   string         `json:"blockHash,omitempty"`
    BlockHeight int64          `json:"blockHeight,omitempty"`
    ExtraInfo   string         `json:"extraInfo,omitempty"`
    MerklePath  *bc.MerklePath `json:"merklePath,omitempty"`
    Timestamp   string         `json:"timestamp,omitempty"`
    TxStatus    TxStatus       `json:"txStatus,omitempty"`
    TxID        string         `json:"txid,omitempty"`
}
```

### PolicyQuote

#### PolicyQuoteResponse

```go
type PolicyQuoteResponse struct {
    Miner     string         `json:"miner"`
    Policy    PolicyResponse `json:"policy"`
    Timestamp string         `json:"timestamp"`
}
```

#### PolicyResponse

```go
type PolicyResponse struct {
    MaxScriptSizePolicy    int64             `json:"maxscriptsizepolicy"`
    MaxTxSigOpsCountPolicy int64             `json:"maxtxsigopscountspolicy"`
    MaxTxSizePolicy        int64             `json:"maxtxsizepolicy"`
    MiningFee              MiningFeeResponse `json:"miningFee"`
}
```

#### MiningFeeResponse

```go
type MiningFeeResponse struct {
    Bytes    int64 `json:"bytes"`
    Satoshis int64 `json:"satoshis"`
}
```

### FeeQuote

#### FeeQuoteResponse

```go
type FeeQuote struct {
    Miner     string            `json:"miner"`
    MiningFee MiningFeeResponse `json:"miningFee"`
    Timestamp string            `json:"timestamp"`
}
```

### SubmitTx

#### SubmitTxResponse

```go
type SubmitTxResponse struct {
    Miner       string         `json:"miner"`
    BlockHash   string         `json:"blockHash,omitempty"`
    BlockHeight int64          `json:"blockHeight,omitempty"`
    ExtraInfo   string         `json:"extraInfo,omitempty"`
    MerklePath  *bc.MerklePath `json:"merklePath,omitempty"`
    Timestamp   string         `json:"timestamp,omitempty"`
    TxStatus    TxStatus       `json:"txStatus,omitempty"`
    TxID        string         `json:"txid,omitempty"`
    Status      int            `json:"status,omitempty"`
    Title       string         `json:"title,omitempty"`
}
```

#### Transaction

  ```go
type Transaction struct {
    CallBackEncryption string   `json:"callBackEncryption,omitempty"`
    CallBackToken      string   `json:"callBackToken,omitempty"`
    CallBackURL        string   `json:"callBackUrl,omitempty"`
    DsCheck            bool     `json:"dsCheck,omitempty"`
    MerkleFormat       string   `json:"merkleFormat,omitempty"`
    MerkleProof        bool     `json:"merkleProof,omitempty"`
    Hex                string   `json:"hex"`
    WaitForStatus      TxStatus `json:"waitForStatus,omitempty"`
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


## MockClientBuilder

Mock Client allows you to test your code without using an actual client and without connecting to any nodes.

### WithMockArc Method

This method allows you to create a client with a different Mock Type passed as parameter.

```go
client := broadcast_client_mock.Builder().
	WithMockArc(broadcast_client_mock.MockSucces).
	Build()
```

| MockType      | Description
|---------------|-----------------------------------------------------------------------------------
| `MockSucces`  | Client will return a successful response from all methods.
| `MockFailure` | Client will return an error that no miner returned a response from all methods.
| `MockTimeout` | Client will return a successful response after a timeout from all methods.

#### MockTimeout

MockTimeout will return a successfull response after around ~10ms more than the timeout provided in the context that is passed to it's method.


Example:

```go
client := broadcast_client_mock.Builder().
	WithMockArc(broadcast_client_mock.MockTimeout).
	Build()

ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // Creating a timeout context with 2 seconds timeout
defer cancel()

result, err := client.GetPolicyQuote(ctx) // client will return a response after around 2 seconds and 10 milliseconds, therefore exceeding the timeout
```

If you pass the context without a timeout, the client will instantly return a successful response (just like from a MockSuccess type).

### Mock Responses Constants

In order to test the Mock Client, you may often need to check if the response from the client is correct. Therefore, some useful constants are exposed from `fixtures` package.

```go
const (
	ProviderMain                                 = "https://mocked_arc_api_url.com/arc"
	ProviderSecondary                            = "https://secondary_mocked_arc_api_url.com/arc"
	MaxScriptSizePolicy       int64              = 100000000
	MaxTxSigOpsCountPolicy    int64              = 4294967295
	MaxTxSizePolicy           int64              = 100000000
	MaxTxSizePolicySecondary  int64              = 220000000
	MiningFeeBytes            int64              = 1000
	SatoshisPerBytes          int64              = 1
	SatoshisPerBytesSecondary int64              = 2
	Timestamp                                    = "2023-09-05T17:03:49.537230128Z"
	TimestampSecondary                           = "2023-09-05T17:05:29.736256927Z"
	TxResponseStatus                             = 200
	TxResponseTitle                              = "OK"
	TxStatus                  broadcast.TxStatus = "SEEN_ON_NETWORK"
	TxBlockHash                                  = "0000000000000000019a575e0ea4d9bbe251dd24c473a0d8407935973151f282"
	TxBlockHashSecondary                         = "0000000000000000045c969f3acd5db37896aba95f91389f2d191496bf15584b"
	TxBlockHeight             int64              = 800182
	TxBlockHeightSecondary    int64              = 799439
	TxExtraInfo                                  = ""
	TxMerklePath                                 = "fd123511fc17127984cea2846b6a6b3965cbd523ea92736c11022372dc922d0013ee15f577bb7d832ce6102e22a15c6e4f4e4214e246d769f04baaf458e3407f5117de8eceb78782a1f90bfaa103915dcc3aa940c7e1476398db47dc62833f3149c51593cf9eeaec36f86c871d64e849a486b97eee1ee863c18aa28e599e3b30853a8563a87441202a97f707d0614cfc024f73ba4a2848ec63e19bc5fb168e994c2c2c105b5951dee4775f20879cb63b3b525370b35200e9f91e9e0fa234dfae5cc4792d17ab08f5cb29aa97e99742dc1e48a9d7e648be33db2d75b806519d88c62a57b39ef900cf68c9d765a1a3b17895d05077ab33280fdb4c09f818d31065b1325b6eb0366c4fd48c2840c1e3a5eb64a9596ffe53bbb315b44784882fb727a98e791d4c8a0bc1d5a88fd85365afdf83da86216cf81b587ee960b5c2a4174dfb26fa2606addf5e51edea58f76cc415f21876c9fa4551891aab537f4e463f1ded910d6a27ff27531c32b4dad15360b35d4d6da348c0c311aaf412bd19a92e2253e90c67104c8f0c0167cf43032f4ffa616cb8c5c28733a3faef5afb83d3c621b252112cbf35859c7b3b4caf820ae43f7054db87da2b7bc735f41b0969b4caadd1cc78be66fb367990103b5682f33605628dcc9af941bf555c64e1c48d3d0955f2c53f7ed7611d976d3effb053eb13b85ca2bff429516aa21227c39d6b6f4106f513069a559805977dee6462d354522d6e42f30d60744bd0f08471577ea5e1f63f319261"
	TxId                                         = "680d975a403fd9ec90f613e87d17802c029d2d930df1c8373cdcdda2f536a1c0"
	TxIdSecondary                                = "d4f4bcb8decdbc75a3dcc54580c76457ee663a28b843956733570b38f7c73b5a"
)
```
