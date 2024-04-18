package mocks

import (
	"context"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

type ArcClientMockFailure struct{}

// GetFeeQuote returns an error.
func (*ArcClientMockFailure) GetFeeQuote(ctx context.Context) ([]*broadcast_api.FeeQuote, *broadcast_api.SubmitFailure) {
	return nil, broadcast.Failure("", broadcast.ErrNoMinerResponse)
}

// GetPolicyQuote returns an error.
func (*ArcClientMockFailure) GetPolicyQuote(ctx context.Context) ([]*broadcast_api.PolicyQuoteResponse, *broadcast_api.SubmitFailure) {
	return nil, broadcast.Failure("", broadcast.ErrNoMinerResponse)
}

// QueryTransaction returns an error.
func (*ArcClientMockFailure) QueryTransaction(ctx context.Context, txID string) (*broadcast_api.QueryTxResponse, *broadcast_api.SubmitFailure) {
	return nil, broadcast.Failure("", broadcast.ErrNoMinerResponse)
}

// SubmitBatchTransactions returns an error.
func (*ArcClientMockFailure) SubmitBatchTransactions(ctx context.Context, tx []*broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitBatchTxResponse, *broadcast_api.SubmitFailure) {
	return nil, broadcast.Failure("", broadcast.ErrAllBroadcastersFailed)
}

// SubmitTransaction returns an error.
func (*ArcClientMockFailure) SubmitTransaction(ctx context.Context, tx *broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitTxResponse, *broadcast_api.SubmitFailure) {
	return nil, broadcast.Failure("", broadcast.ErrAllBroadcastersFailed)
}

// NewArcClientMockFailure creates a new mock arc client that returns an error from all its methods.
func NewArcClientMockFailure() broadcast_api.Client {
	return &ArcClientMockFailure{}
}
