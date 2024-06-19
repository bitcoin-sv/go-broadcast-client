package mocks

import (
	"context"

	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock/fixtures"
)

type ArcClientMockNilQueryTxResp struct{}

// GetFeeQuote returns a successful FeeQuote response.
func (*ArcClientMockNilQueryTxResp) GetFeeQuote(ctx context.Context) ([]*broadcast_api.FeeQuote, error) {
	quotes := []*broadcast_api.FeeQuote{Fee1, Fee2}
	return quotes, nil
}

// GetPolicyQuote return a successful PolicyQuoteResponse.
func (*ArcClientMockNilQueryTxResp) GetPolicyQuote(ctx context.Context) ([]*broadcast_api.PolicyQuoteResponse, error) {
	policies := []*broadcast_api.PolicyQuoteResponse{Policy1, Policy2}
	return policies, nil
}

// QueryTransaction returns a successful QueryTxResponse.
func (*ArcClientMockNilQueryTxResp) QueryTransaction(ctx context.Context, txID string) (*broadcast_api.QueryTxResponse, error) {
	return nil, nil
}

// SubmitTransaction returns a successful SubmitTxResponse.
func (*ArcClientMockNilQueryTxResp) SubmitTransaction(ctx context.Context, tx *broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitTxResponse, error) {
	return &broadcast_api.SubmitTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: fixtures.ProviderMain},
		SubmittedTx:  SubmittedTx,
	}, nil
}

// SubmitBatchTransactions returns a successful SubmitBatchTxResponse.
func (*ArcClientMockNilQueryTxResp) SubmitBatchTransactions(ctx context.Context, tx []*broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitBatchTxResponse, error) {
	return &broadcast_api.SubmitBatchTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: fixtures.ProviderMain},
		Transactions: []*broadcast_api.SubmittedTx{
			SubmittedTx,
			SubmittedTxSecondary,
		},
	}, nil
}

func NewArcClientMockNilQueryTxResp() broadcast_api.Client {
	return &ArcClientMockNilQueryTxResp{}
}
