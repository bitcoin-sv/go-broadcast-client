package mocks

import (
	"context"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

type ArcClientMockFailure struct{}

// GetFeeQuote implements broadcast.Client.
func (*ArcClientMockFailure) GetFeeQuote(ctx context.Context) ([]*broadcast_api.FeeQuote, error) {
	return nil, broadcast.ErrNoMinerResponse
}

// GetPolicyQuote implements broadcast.Client.
func (*ArcClientMockFailure) GetPolicyQuote(ctx context.Context) ([]*broadcast_api.PolicyQuoteResponse, error) {
	return nil, broadcast.ErrNoMinerResponse
}

// QueryTransaction implements broadcast.Client.
func (*ArcClientMockFailure) QueryTransaction(ctx context.Context, txID string) (*broadcast_api.QueryTxResponse, error) {
	return nil, broadcast.ErrAllBroadcastersFailed
}

// SubmitBatchTransactions implements broadcast.Client.
func (*ArcClientMockFailure) SubmitBatchTransactions(ctx context.Context, tx []*broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitBatchTxResponse, error) {
	return nil, broadcast.ErrAllBroadcastersFailed
}

// SubmitTransaction implements broadcast.Client.
func (*ArcClientMockFailure) SubmitTransaction(ctx context.Context, tx *broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitTxResponse, error) {
	return nil, broadcast.ErrAllBroadcastersFailed
}

func NewArcClientMockFailure() broadcast_api.Client {
	return &ArcClientMockFailure{}
}
