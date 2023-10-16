package mocks

import (
	"github.com/libsv/go-bc"

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

	mp, _  = bc.NewMerklePathFromStr(fixtures.TxMerklePath)
	mp2, _ = bc.NewMerklePathFromStr(fixtures.TxMerklePathSecondary)

	SubmittedTx = &broadcast_api.SubmittedTx{
		BaseSubmitTxResponse: broadcast_api.BaseSubmitTxResponse{
			Status: fixtures.TxResponseStatus,
			Title:  fixtures.TxResponseTitle,
			BaseTxResponse: broadcast_api.BaseTxResponse{
				BlockHash:   fixtures.TxBlockHash,
				BlockHeight: fixtures.TxBlockHeight,
				ExtraInfo:   fixtures.TxExtraInfo,
				MerklePath:  fixtures.TxMerklePath,
				Timestamp:   fixtures.Timestamp,
				TxStatus:    fixtures.TxStatus,
				TxID:        fixtures.TxId,
			},
		},
		MerklePath: mp,
	}

	SubmittedTxSecondary = &broadcast_api.SubmittedTx{
		BaseSubmitTxResponse: broadcast_api.BaseSubmitTxResponse{
			Status: fixtures.TxResponseStatus,
			Title:  fixtures.TxResponseTitle,
			BaseTxResponse: broadcast_api.BaseTxResponse{
				BlockHash:   fixtures.TxBlockHashSecondary,
				BlockHeight: fixtures.TxBlockHeightSecondary,
				ExtraInfo:   fixtures.TxExtraInfo,
				MerklePath:  fixtures.TxMerklePath,
				Timestamp:   fixtures.Timestamp,
				TxStatus:    fixtures.TxStatus,
				TxID:        fixtures.TxIdSecondary,
			},
		},
		MerklePath: mp2,
	}
)

func QueryTx(txID string) *broadcast_api.QueryTxResponse {
	return &broadcast_api.QueryTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: fixtures.ProviderMain},
		BaseTxResponse: broadcast_api.BaseTxResponse{
			BlockHash:   fixtures.TxBlockHash,
			BlockHeight: fixtures.TxBlockHeight,
			ExtraInfo:   fixtures.TxExtraInfo,
			MerklePath:  fixtures.TxMerklePath,
			Timestamp:   fixtures.Timestamp,
			TxStatus:    fixtures.TxStatus,
			TxID:        txID,
		},
	}
}
