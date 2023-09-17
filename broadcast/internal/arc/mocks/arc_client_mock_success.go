package mocks

import (
	"context"

	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

type ArcClientMock struct{}

// GetFeeQuote returns a successful FeeQuote response.
func (*ArcClientMock) GetFeeQuote(ctx context.Context) ([]*broadcast_api.FeeQuote, error) {
	quote1 := &broadcast_api.FeeQuote{
		BaseResponse: broadcast_api.BaseResponse{Miner: MockedApiUrl1},
		MiningFee:    Policy1.Policy.MiningFee,
		Timestamp:    Policy1.Timestamp,
	}

	quote2 := &broadcast_api.FeeQuote{
		BaseResponse: broadcast_api.BaseResponse{Miner: MockedApiUrl2},
		MiningFee:    Policy2.Policy.MiningFee,
		Timestamp:    Policy2.Timestamp,
	}

	quotes := make([]*broadcast_api.FeeQuote, 2)
	quotes = append(quotes, quote1)
	quotes = append(quotes, quote2)

	return quotes, nil
}

// GetPolicyQuote return a successful PolicyQuoteResponse.
func (*ArcClientMock) GetPolicyQuote(ctx context.Context) ([]*broadcast_api.PolicyQuoteResponse, error) {
	policies := make([]*broadcast_api.PolicyQuoteResponse, 2)
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
		BaseResponse: broadcast_api.BaseResponse{Miner: MockedApiUrl1},
		SubmittedTx:  SubmittedTx,
	}, nil
}

// SubmitBatchTransactions returns a successful SubmitBatchTxResponse.
func (*ArcClientMock) SubmitBatchTransactions(ctx context.Context, tx []*broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitBatchTxResponse, error) {
	return &broadcast_api.SubmitBatchTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: MockedApiUrl1},
		Transactions: []*broadcast_api.SubmittedTx{
			SubmittedTx,
			SubmittedTx,
		},
	}, nil
}

func NewArcClientMock() broadcast_api.Client {
	return &ArcClientMock{}
}
