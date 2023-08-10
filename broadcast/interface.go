package broadcast

import (
	"context"
	"time"
)

type BestQuoter interface {
	GetBestQuote(ctx context.Context) (*FeeQuote, error)
}

type FastestQuoter interface {
	GetFastestQuote(ctx context.Context, timeout time.Duration) (*FeeQuote, error)
}

type FeeQuoter interface {
	GetFeeQuote(ctx context.Context) (*FeeQuote, error)
}

type PolicyQuoter interface {
	GetPolicyQuote(ctx context.Context) (*PolicyQuoteResponse, error)
}

type TransactionQuerier interface {
	QueryTransaction(ctx context.Context, txID string) (*QueryTxResponse, error)
}

type TransactionSubmitter interface {
	SubmitTransaction(ctx context.Context, tx *Transaction) (*SubmitTxResponse, error)
}

type TransactionsSubmitter interface {
	SubmitBatchTransactions(ctx context.Context, tx []*Transaction) ([]*SubmitTxResponse, error)
}

type Client interface {
	BestQuoter
	FastestQuoter
	FeeQuoter
	PolicyQuoter
	TransactionQuerier
	TransactionSubmitter
	TransactionsSubmitter
}
