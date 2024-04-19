package mocks

import (
	"context"

	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

type ArcClientMockFailure struct{}

// GetFeeQuote returns an error.
func (*ArcClientMockFailure) GetFeeQuote(ctx context.Context) ([]*broadcast_api.FeeQuote, *broadcast_api.FailureResponse) {
	return nil, broadcast_api.Failure("", broadcast_api.ErrNoMinerResponse)
}

// GetPolicyQuote returns an error.
func (*ArcClientMockFailure) GetPolicyQuote(ctx context.Context) ([]*broadcast_api.PolicyQuoteResponse, *broadcast_api.FailureResponse) {
	return nil, broadcast_api.Failure("", broadcast_api.ErrNoMinerResponse)
}

// QueryTransaction returns an error.
func (*ArcClientMockFailure) QueryTransaction(ctx context.Context, txID string) (*broadcast_api.QueryTxResponse, *broadcast_api.FailureResponse) {
	return nil, broadcast_api.Failure("", broadcast_api.ErrNoMinerResponse)
}

// SubmitBatchTransactions returns an error.
func (*ArcClientMockFailure) SubmitBatchTransactions(ctx context.Context, tx []*broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitBatchTxResponse, *broadcast_api.FailureResponse) {
	return nil, broadcast_api.Failure("", broadcast_api.ErrAllBroadcastersFailed)
}

// SubmitTransaction returns an error.
func (*ArcClientMockFailure) SubmitTransaction(ctx context.Context, tx *broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitTxResponse, *broadcast_api.FailureResponse) {
	return nil, broadcast_api.Failure("", broadcast_api.ErrAllBroadcastersFailed)
}

// NewArcClientMockFailure creates a new mock arc client that returns an error from all its methods.
func NewArcClientMockFailure() broadcast_api.Client {
	return &ArcClientMockFailure{}
}
