package models

import "github.com/bitcoin-sv/go-broadcast-client/common"

type QueryTxResponse struct {
	BlockHash   string          `json:"blockHash,omitempty"`
	BlockHeight int64           `json:"blockHeight,omitempty"`
	Timestamp   string          `json:"timestamp,omitempty"`
	TxID        string          `json:"txid,omitempty"`
	TxStatus    common.TxStatus `json:"txStatus,omitempty"`
}
