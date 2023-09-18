package mocks

import (
	"context"
	"time"

	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

type ArcClientMockTimeout struct{}

// GetFeeQuote returns a successful FeeQuote response.
func (*ArcClientMockTimeout) GetFeeQuote(ctx context.Context) ([]*broadcast_api.FeeQuote, error) {
	if deadline, ok := ctx.Deadline(); ok {
		time.Sleep(time.Until(deadline) + 10*time.Millisecond)
	}

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
func (*ArcClientMockTimeout) GetPolicyQuote(ctx context.Context) ([]*broadcast_api.PolicyQuoteResponse, error) {
	if deadline, ok := ctx.Deadline(); ok {
		time.Sleep(time.Until(deadline) + 10*time.Millisecond)
	}

	policies := make([]*broadcast_api.PolicyQuoteResponse, 2)
	policies = append(policies, Policy1)
	policies = append(policies, Policy2)

	return policies, nil
}

// QueryTransaction returns a successful QueryTxResponse.
func (*ArcClientMockTimeout) QueryTransaction(ctx context.Context, txID string) (*broadcast_api.QueryTxResponse, error) {
	if deadline, ok := ctx.Deadline(); ok {
		time.Sleep(time.Until(deadline) + 10*time.Millisecond)
	}

	return QueryTx(txID), nil
}

// SubmitTransaction returns a successful SubmitTxResponse.
func (*ArcClientMockTimeout) SubmitTransaction(ctx context.Context, tx *broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitTxResponse, error) {
	if deadline, ok := ctx.Deadline(); ok {
		time.Sleep(time.Until(deadline) + 10*time.Millisecond)
	}

	return &broadcast_api.SubmitTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: MockedApiUrl1},
		SubmittedTx:  SubmittedTx,
	}, nil
}

// SubmitBatchTransactions returns a successful SubmitBatchTxResponse.
func (*ArcClientMockTimeout) SubmitBatchTransactions(ctx context.Context, tx []*broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitBatchTxResponse, error) {
	if deadline, ok := ctx.Deadline(); ok {
		time.Sleep(time.Until(deadline) + 10*time.Millisecond)
	}

	return &broadcast_api.SubmitBatchTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: MockedApiUrl1},
		Transactions: []*broadcast_api.SubmittedTx{
			SubmittedTx,
			SubmittedTx,
		},
	}, nil
}

func NewArcClientMockTimeout() broadcast_api.Client {
	return &ArcClientMockTimeout{}
}
