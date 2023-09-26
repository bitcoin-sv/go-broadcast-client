package mocks

import (
	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock/fixtures"
)

var (
	policy1 = &broadcast_api.PolicyQuoteResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: fixtures.ProviderMain},
		Policy: broadcast_api.PolicyResponse{
			MaxScriptSizePolicy:    fixtures.MaxScriptSizePolicy,
			MaxTxSigOpsCountPolicy: fixtures.MaxTxSigOpsCountPolicy,
			MaxTxSizePolicy:        fixtures.MaxTxSizePolicy,
			MiningFee: broadcast_api.MiningFeeResponse{
				Bytes:    fixtures.MiningFeeBytes,
				Satoshis: fixtures.SatoshisPerBytes,
			},
		},
		Timestamp: fixtures.Timestamp,
	}

	policy2 = &broadcast_api.PolicyQuoteResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: fixtures.ProviderSecondary},
		Policy: broadcast_api.PolicyResponse{
			MaxScriptSizePolicy:    fixtures.MaxScriptSizePolicy,
			MaxTxSigOpsCountPolicy: fixtures.MaxTxSigOpsCountPolicy,
			MaxTxSizePolicy:        fixtures.MaxTxSizePolicySecondary,
			MiningFee: broadcast_api.MiningFeeResponse{
				Bytes:    fixtures.MiningFeeBytes,
				Satoshis: fixtures.SatoshisPerBytesSecondary,
			},
		},
		Timestamp: fixtures.TimestampSecondary,
	}

	submittedTx = &broadcast_api.SubmittedTx{
		Status:      fixtures.TxResponseStatus,
		Title:       fixtures.TxResponseTitle,
		TxStatus:    fixtures.TxStatus,
		BlockHash:   fixtures.TxBlockHash,
		BlockHeight: fixtures.TxBlockHeight,
		ExtraInfo:   fixtures.TxExtraInfo,
	}

	submittedTxSecondary = &broadcast_api.SubmittedTx{
		Status:      fixtures.TxResponseStatus,
		Title:       fixtures.TxResponseTitle,
		TxStatus:    fixtures.TxStatus,
		BlockHash:   fixtures.TxBlockHashSecondary,
		BlockHeight: fixtures.TxBlockHeightSecondary,
		ExtraInfo:   fixtures.TxExtraInfo,
	}
)

func queryTx(txID string) *broadcast_api.QueryTxResponse {
	return &broadcast_api.QueryTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: fixtures.ProviderMain},
		Timestamp:    fixtures.Timestamp,
		TxStatus:     fixtures.TxStatus,
		BlockHash:    fixtures.TxBlockHash,
		BlockHeight:  fixtures.TxBlockHeight,
		TxID:         txID,
	}
}
