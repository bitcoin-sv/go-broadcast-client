package mocks

import (
	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock/fixtures"
)

var (
	Policy1 = &broadcast_api.PolicyQuoteResponse{
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

	Policy2 = &broadcast_api.PolicyQuoteResponse{
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

	Fee1 = &broadcast_api.FeeQuote{
		BaseResponse: broadcast_api.BaseResponse{Miner: fixtures.ProviderMain},
		MiningFee:    Policy1.Policy.MiningFee,
		Timestamp:    Policy1.Timestamp,
	}

	Fee2 = &broadcast_api.FeeQuote{
		BaseResponse: broadcast_api.BaseResponse{Miner: fixtures.ProviderSecondary},
		MiningFee:    Policy2.Policy.MiningFee,
		Timestamp:    Policy2.Timestamp,
	}

	SubmittedTx = &broadcast_api.SubmittedTx{
		Status:      fixtures.TxResponseStatus,
		Title:       fixtures.TxResponseTitle,
		TxStatus:    fixtures.TxStatus,
		BlockHash:   fixtures.TxBlockHash,
		BlockHeight: fixtures.TxBlockHeight,
		ExtraInfo:   fixtures.TxExtraInfo,
	}

	SubmittedTxSecondary = &broadcast_api.SubmittedTx{
		Status:      fixtures.TxResponseStatus,
		Title:       fixtures.TxResponseTitle,
		TxStatus:    fixtures.TxStatus,
		BlockHash:   fixtures.TxBlockHashSecondary,
		BlockHeight: fixtures.TxBlockHeightSecondary,
		ExtraInfo:   fixtures.TxExtraInfo,
	}
)

func QueryTx(txID string) *broadcast_api.QueryTxResponse {
	return &broadcast_api.QueryTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: fixtures.ProviderMain},
		Timestamp:    fixtures.Timestamp,
		TxStatus:     fixtures.TxStatus,
		BlockHash:    fixtures.TxBlockHash,
		BlockHeight:  fixtures.TxBlockHeight,
		TxID:         txID,
	}
}
