package mocks

import (
	"context"
	"time"

	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

const defaultMockTimeout = 20 * time.Second

type ArcClientMockTimeout struct{}

// GetFeeQuote returns a successful FeeQuote response after a long time.
func (*ArcClientMockTimeout) GetFeeQuote(ctx context.Context) ([]*broadcast_api.FeeQuote, error) {
	quote1 := &broadcast_api.FeeQuote{
		BaseResponse: broadcast_api.BaseResponse{Miner: MockedApiUrl1},
		MiningFee: broadcast_api.MiningFeeResponse{
			Bytes:    1000,
			Satoshis: 1,
		},
		Timestamp: "2023-09-05T17:03:49.537230128Z",
	}

	quote2 := &broadcast_api.FeeQuote{
		BaseResponse: broadcast_api.BaseResponse{Miner: MockedApiUrl2},
		MiningFee: broadcast_api.MiningFeeResponse{
			Bytes:    1000,
			Satoshis: 2,
		},
		Timestamp: "2023-09-05T17:05:29.736256927Z",
	}

	quotes := make([]*broadcast_api.FeeQuote, 2)
	quotes = append(quotes, quote1)
	quotes = append(quotes, quote2)

	time.Sleep(defaultMockTimeout)
	return quotes, nil
}

// GetPolicyQuote return a successful PolicyQuoteResponse after a long time.
func (*ArcClientMockTimeout) GetPolicyQuote(ctx context.Context) ([]*broadcast_api.PolicyQuoteResponse, error) {
	policy1 := &broadcast_api.PolicyQuoteResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: MockedApiUrl1},
		Policy: broadcast_api.PolicyResponse{
			MaxScriptSizePolicy:    100000000,
			MaxTxSigOpsCountPolicy: 4294967295,
			MaxTxSizePolicy:        100000000,
			MiningFee: broadcast_api.MiningFeeResponse{
				Bytes:    1000,
				Satoshis: 1,
			},
		},
		Timestamp: "2023-09-05T17:03:49.537230128Z",
	}

	policy2 := &broadcast_api.PolicyQuoteResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: MockedApiUrl2},
		Policy: broadcast_api.PolicyResponse{
			MaxScriptSizePolicy:    100000000,
			MaxTxSigOpsCountPolicy: 4294967295,
			MaxTxSizePolicy:        220000000,
			MiningFee: broadcast_api.MiningFeeResponse{
				Bytes:    1000,
				Satoshis: 2,
			},
		},
		Timestamp: "2023-09-05T17:05:29.736256927Z",
	}

	policies := make([]*broadcast_api.PolicyQuoteResponse, 2)
	policies = append(policies, policy1)
	policies = append(policies, policy2)

	time.Sleep(defaultMockTimeout)
	return policies, nil
}

// QueryTransaction returns a successful QueryTxResponse after a long time.
func (*ArcClientMockTimeout) QueryTransaction(ctx context.Context, txID string) (*broadcast_api.QueryTxResponse, error) {
	return &broadcast_api.QueryTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: MockedApiUrl1},
		Timestamp:    "2023-09-05T17:05:29.736256927Z",
		TxID:         txID,
		TxStatus:     "SEEN_ON_NETWORK",
	}, nil
}

// SubmitTransaction returns a successful SubmitTxResponse after a long time.
func (*ArcClientMockTimeout) SubmitTransaction(ctx context.Context, tx *broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitTxResponse, error) {
	time.Sleep(defaultMockTimeout)
	return &broadcast_api.SubmitTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: MockedApiUrl1},
		SubmittedTx: &broadcast_api.SubmittedTx{
			Status:   200,
			Title:    "OK",
			TxStatus: "SENT_TO_NETWORK",
		},
	}, nil
}

// SubmitBatchTransactions returns a successful SubmitBatchTxResponse after a long time.
func (*ArcClientMockTimeout) SubmitBatchTransactions(ctx context.Context, tx []*broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitBatchTxResponse, error) {
	time.Sleep(defaultMockTimeout)
	return &broadcast_api.SubmitBatchTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: MockedApiUrl1},
		Transactions: []*broadcast_api.SubmittedTx{
			{
				Status:   200,
				Title:    "OK",
				TxStatus: "SENT_TO_NETWORK",
			},
			{
				Status:   200,
				Title:    "OK",
				TxStatus: "SENT_TO_NETWORK",
			},
		},
	}, nil
}

func NewArcClientMockTimeout() broadcast_api.Client {
	return &ArcClientMockTimeout{}
}
