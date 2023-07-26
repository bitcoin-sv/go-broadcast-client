package broadcast

import (
	"context"

	"github.com/bitcoin-sv/go-broadcast-client/models"
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
	QueryTransaction(ctx context.Context, txID string) (*models.QueryTxResponse, error)
}

type TransactionSubmitter interface {
}

type TransactionsSubmitter interface {
}
