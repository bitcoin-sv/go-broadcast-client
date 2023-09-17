package mocks

import broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"

const (
	MockedApiUrl1 = "https://mocked_api_url.com/arc"
	MockedApiUrl2 = "https://second_mocked_api_url.com/arc"
)

var (
	Policy1 = &broadcast_api.PolicyQuoteResponse{
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

	Policy2 = &broadcast_api.PolicyQuoteResponse{
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

	SubmittedTx = &broadcast_api.SubmittedTx{
		Status:   200,
		Title:    "OK",
		TxStatus: "SENT_TO_NETWORK",
	}
)

func QueryTx(txID string) *broadcast_api.QueryTxResponse {
	return &broadcast_api.QueryTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: MockedApiUrl1},
		Timestamp:    "2023-09-05T17:05:29.736256927Z",
		TxID:         txID,
		TxStatus:     "SEEN_ON_NETWORK",
	}
}
