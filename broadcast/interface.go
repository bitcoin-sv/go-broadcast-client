package broadcast

import (
	"context"
)

type FeeQuoter interface {
	GetFeeQuote(ctx context.Context) ([]*FeeQuote, error)
}

type PolicyQuoter interface {
	GetPolicyQuote(ctx context.Context) ([]*PolicyQuoteResponse, error)
}

// TransactionQuerier is the interface that wraps the QueryTransaction method.
// It takes a transaction ID and returns the transaction details, like it's status, hash, height etc.
// Everything is wrapped in the QueryTxResponse struct.
type TransactionQuerier interface {
	QueryTransaction(ctx context.Context, txID string) (*QueryTxResponse, error)
}

// TransactionSubmitter is the interface that wraps the SubmitTransaction method.
// It takes a transaction and tries to broadcast it to the P2P network.
// Transaction object needs RawTx to be set. All other fields are optional and used to append headers related to status callbacks.
// As a result it returns a SubmitTxResponse object.
type TransactionSubmitter interface {
	SubmitTransaction(ctx context.Context, tx *Transaction, opts ...TransactionOptFunc) (*SubmitTxResponse, error)
}

// TransactionsSubmitter is the interface that wraps the SubmitBatchTransactions method.
// It is the same as TransactionSubmitter but it takes a slice of transactions and tries to broadcast them to the P2P network.
// As a result it returns a slice of SubmitTxResponse objects.
type TransactionsSubmitter interface {
	SubmitBatchTransactions(ctx context.Context, tx []*Transaction, opts ...TransactionOptFunc) ([]*SubmitTxResponse, error)
}

// Client is a grouping interface that represents the entire exposed functionality of the broadcast client.
type Client interface {
	FeeQuoter
	PolicyQuoter
	TransactionQuerier
	TransactionSubmitter
	TransactionsSubmitter
}
