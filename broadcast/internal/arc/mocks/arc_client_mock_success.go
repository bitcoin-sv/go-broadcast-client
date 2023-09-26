package mocks

import (
	"context"

	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock/fixtures"
)

type ArcClientMock struct{}

// GetFeeQuote returns a successful FeeQuote response.
func (*ArcClientMock) GetFeeQuote(ctx context.Context) ([]*broadcast_api.FeeQuote, error) {
	quotes := make([]*broadcast_api.FeeQuote, 0)
	quotes = append(quotes, Fee1)
	quotes = append(quotes, Fee2)

	return quotes, nil
}

// GetPolicyQuote return a successful PolicyQuoteResponse.
func (*ArcClientMock) GetPolicyQuote(ctx context.Context) ([]*broadcast_api.PolicyQuoteResponse, error) {
	policies := make([]*broadcast_api.PolicyQuoteResponse, 0)
	policies = append(policies, Policy1)
	policies = append(policies, Policy2)

	return policies, nil
}

// QueryTransaction returns a successful QueryTxResponse.
func (*ArcClientMock) QueryTransaction(ctx context.Context, txID string) (*broadcast_api.QueryTxResponse, error) {
	return QueryTx(txID), nil
}

// SubmitTransaction returns a successful SubmitTxResponse.
func (*ArcClientMock) SubmitTransaction(ctx context.Context, tx *broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitTxResponse, error) {
	return &broadcast_api.SubmitTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: fixtures.ProviderMain},
		SubmittedTx:  SubmittedTx,
	}, nil
}

// SubmitBatchTransactions returns a successful SubmitBatchTxResponse.
func (*ArcClientMock) SubmitBatchTransactions(ctx context.Context, tx []*broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitBatchTxResponse, error) {
	return &broadcast_api.SubmitBatchTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: fixtures.ProviderMain},
		Transactions: []*broadcast_api.SubmittedTx{
			SubmittedTx,
			SubmittedTxSecondary,
		},
	}, nil
}

func NewArcClientMock() broadcast_api.Client {
	return &ArcClientMock{}
}
