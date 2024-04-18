package mocks

import (
	"context"
	"time"

	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock/fixtures"
)

type ArcClientMockTimeout struct{}

// GetFeeQuote returns a successful FeeQuote response.
func (*ArcClientMockTimeout) GetFeeQuote(ctx context.Context) ([]*broadcast_api.FeeQuote, *broadcast_api.SubmitFailure) {
	if deadline, ok := ctx.Deadline(); ok {
		time.Sleep(time.Until(deadline) + 10*time.Millisecond)
	}

	quotes := make([]*broadcast_api.FeeQuote, 0)
	quotes = append(quotes, Fee1)
	quotes = append(quotes, Fee2)

	return quotes, nil
}

// GetPolicyQuote return a successful PolicyQuoteResponse.
func (*ArcClientMockTimeout) GetPolicyQuote(ctx context.Context) ([]*broadcast_api.PolicyQuoteResponse, *broadcast_api.SubmitFailure) {
	if deadline, ok := ctx.Deadline(); ok {
		time.Sleep(time.Until(deadline) + 10*time.Millisecond)
	}

	policies := make([]*broadcast_api.PolicyQuoteResponse, 0)
	policies = append(policies, Policy1)
	policies = append(policies, Policy2)

	return policies, nil
}

// QueryTransaction returns a successful QueryTxResponse.
func (*ArcClientMockTimeout) QueryTransaction(ctx context.Context, txID string) (*broadcast_api.QueryTxResponse, *broadcast_api.SubmitFailure) {
	if deadline, ok := ctx.Deadline(); ok {
		time.Sleep(time.Until(deadline) + 10*time.Millisecond)
	}

	return QueryTx(txID), nil
}

// SubmitTransaction returns a successful SubmitTxResponse.
func (*ArcClientMockTimeout) SubmitTransaction(ctx context.Context, tx *broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitTxResponse, *broadcast_api.SubmitFailure) {
	if deadline, ok := ctx.Deadline(); ok {
		time.Sleep(time.Until(deadline) + 10*time.Millisecond)
	}

	return &broadcast_api.SubmitTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: fixtures.ProviderMain},
		SubmittedTx:  SubmittedTx,
	}, nil
}

// SubmitBatchTransactions returns a successful SubmitBatchTxResponse.
func (*ArcClientMockTimeout) SubmitBatchTransactions(ctx context.Context, tx []*broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitBatchTxResponse, *broadcast_api.SubmitFailure) {
	if deadline, ok := ctx.Deadline(); ok {
		time.Sleep(time.Until(deadline) + 10*time.Millisecond)
	}

	return &broadcast_api.SubmitBatchTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: fixtures.ProviderMain},
		Transactions: []*broadcast_api.SubmittedTx{
			SubmittedTx,
			SubmittedTxSecondary,
		},
	}, nil
}

func NewArcClientMockTimeout() broadcast_api.Client {
	return &ArcClientMockTimeout{}
}
