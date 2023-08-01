package broadcast

type QueryTxResponse struct {
	BlockHash   string   `json:"blockHash,omitempty"`
	BlockHeight int64    `json:"blockHeight,omitempty"`
	Timestamp   string   `json:"timestamp,omitempty"`
	TxID        string   `json:"txid,omitempty"`
	TxStatus    TxStatus `json:"txStatus,omitempty"`
}
