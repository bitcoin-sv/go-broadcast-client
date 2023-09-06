package mocks

import (
	"context"

	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

type ArcClientMockOnChain struct {
	ArcClientMock
}

func (*ArcClientMockOnChain) SubmitTransaction(ctx context.Context, tx *broadcast_api.Transaction, opts ...broadcast_api.TransactionOptFunc) (*broadcast_api.SubmitTxResponse, error) {
	return &broadcast_api.SubmitTxResponse{
		BaseResponse: broadcast_api.BaseResponse{Miner: MockedApiUrl1},
		SubmittedTx: &broadcast_api.SubmittedTx{
			Status:   200,
			Title:    "OK",
			TxStatus: "CONFIRMED",
		},
	}, nil
}

func NewArcClientMockOnChain() broadcast_api.Client {
	return &ArcClientMockOnChain{}
}
