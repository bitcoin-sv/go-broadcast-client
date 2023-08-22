package broadcast

// QueryTxResponse is the response returned by the QueryTransaction method.
type QueryTxResponse struct {
	// Miner is the URL of the miner that returned the transaction.
	Miner string `json:"miner"`
	// BlockHash is the hash of the block where the transaction was included.
	BlockHash string `json:"blockHash,omitempty"`
	// BlockHeight is the height of the block where the transaction was included.
	BlockHeight int64 `json:"blockHeight,omitempty"`
	// Timestamp is the timestamp of the block where the transaction was included.
	Timestamp string `json:"timestamp,omitempty"`
	// TxID is the transaction id.
	TxID string `json:"txid,omitempty"`
	// TxStatus is the status of the transaction.
	TxStatus TxStatus `json:"txStatus,omitempty"`
}
