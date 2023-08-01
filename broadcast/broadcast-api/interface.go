package broadcast_api

import (
	"context"
)

type BestQuoter interface {
}

type FastestQuoter interface {
}

type FeeQuoter interface {
}

type PolicyQuoter interface {
}

type TransactionQuerier interface {
	// Think about adding TransactionQueryOpts here if clients implement handling it in future
	QueryTransaction(ctx context.Context, txID string) (*QueryTxResponse, error)
}

type TransactionSubmitter interface {
	SubmitTransaction(ctx context.Context, tx *Transaction) (*SubmitTxResponse, error)
}

type TransactionsSubmitter interface {
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
