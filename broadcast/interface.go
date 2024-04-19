package broadcast

import (
	"context"
)

// FeeQuoter it the interface that wraps GetFeeQuote method.
// It retrieves the Fee Quote from the configured miners.
type FeeQuoter interface {
	GetFeeQuote(ctx context.Context) ([]*FeeQuote, *FailureResponse)
}

// PolicyQuoter it the interface that wraps GetPolicyQuote method.
// It retrieves the Policy Quote from the configured miners.
type PolicyQuoter interface {
	GetPolicyQuote(ctx context.Context) ([]*PolicyQuoteResponse, *FailureResponse)
}

// TransactionQuerier is the interface that wraps the QueryTransaction method.
// It takes a transaction ID and returns the transaction details, like it's status, hash, height etc.
// Everything is wrapped in the QueryTxResponse struct.
type TransactionQuerier interface {
	QueryTransaction(ctx context.Context, txID string) (*QueryTxResponse, *FailureResponse)
}

// TransactionSubmitter is the interface that wraps the SubmitTransaction method.
// It takes a transaction and tries to broadcast it to the P2P network.
// Transaction object needs RawTx to be set. All other fields are optional and used to append headers related to status callbacks.
// As a result it returns a SubmitTxResponse object.
type TransactionSubmitter interface {
	SubmitTransaction(ctx context.Context, tx *Transaction, opts ...TransactionOptFunc) (*SubmitTxResponse, *FailureResponse)
}

// TransactionsSubmitter is the interface that wraps the SubmitBatchTransactions method.
// It is the same as TransactionSubmitter but it takes a slice of transactions and tries to broadcast them to the P2P network.
// As a result it returns a SubmitBatchTxResponse, which includes a slice of SubmitTxResponse objects.
type TransactionsSubmitter interface {
	SubmitBatchTransactions(ctx context.Context, tx []*Transaction, opts ...TransactionOptFunc) (*SubmitBatchTxResponse, *FailureResponse)
}

// Client is a grouping interface that represents the entire exposed functionality of the broadcast client.
type Client interface {
	FeeQuoter
	PolicyQuoter
	TransactionQuerier
	TransactionSubmitter
	TransactionsSubmitter
}
