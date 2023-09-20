package mocks

import (
	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	mock_consts "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock/mock-consts"
)

var (
	policy1 = &broadcast_api.PolicyQuoteResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: mock_consts.MockedProviderMain},
		Policy: broadcast_api.PolicyResponse{
			MaxScriptSizePolicy:    mock_consts.MockedMaxScriptSizePolicy,
			MaxTxSigOpsCountPolicy: mock_consts.MockedMaxTxSigOpsCountPolicy,
			MaxTxSizePolicy:        mock_consts.MockedMaxTxSizePolicy,
			MiningFee: broadcast_api.MiningFeeResponse{
				Bytes:    mock_consts.MockedMiningFeeBytes,
				Satoshis: mock_consts.MockedSatoshisPerBytes,
			},
		},
		Timestamp: mock_consts.MockedTimestamp,
	}

	policy2 = &broadcast_api.PolicyQuoteResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: mock_consts.MockedProviderSecondary},
		Policy: broadcast_api.PolicyResponse{
			MaxScriptSizePolicy:    mock_consts.MockedMaxScriptSizePolicy,
			MaxTxSigOpsCountPolicy: mock_consts.MockedMaxTxSigOpsCountPolicy,
			MaxTxSizePolicy:        mock_consts.MockedMaxTxSizePolicySecondary,
			MiningFee: broadcast_api.MiningFeeResponse{
				Bytes:    mock_consts.MockedMiningFeeBytes,
				Satoshis: mock_consts.MockedSatoshisPerBytesSecondary,
			},
		},
		Timestamp: mock_consts.MockedTimestampSecondary,
	}

	submittedTx = &broadcast_api.SubmittedTx{
		Status:      mock_consts.MockedTxResponseStatus,
		Title:       mock_consts.MockedTxResponseTitle,
		TxStatus:    mock_consts.MockedTxStatus,
		BlockHash:   mock_consts.MockedTxBlockHash,
		BlockHeight: mock_consts.MockedTxBlockHeight,
		ExtraInfo:   mock_consts.MockedTxExtraInfo,
	}

	submittedTxSecondary = &broadcast_api.SubmittedTx{
		Status:      mock_consts.MockedTxResponseStatus,
		Title:       mock_consts.MockedTxResponseTitle,
		TxStatus:    mock_consts.MockedTxStatus,
		BlockHash:   mock_consts.MockedTxBlockHashSecondary,
		BlockHeight: mock_consts.MockedTxBlockHeightSecondary,
		ExtraInfo:   mock_consts.MockedTxExtraInfo,
	}
)

func queryTx(txID string) *broadcast_api.QueryTxResponse {
	return &broadcast_api.QueryTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: mock_consts.MockedProviderMain},
		Timestamp:    mock_consts.MockedTimestamp,
		TxStatus:     mock_consts.MockedTxStatus,
		BlockHash:    mock_consts.MockedTxBlockHash,
		BlockHeight:  mock_consts.MockedTxBlockHeight,
		TxID:         txID,
	}
}
