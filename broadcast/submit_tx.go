package broadcast

// SubmitTxResponse is the response returned by the SubmitTransaction method.
type SubmitTxResponse struct {
	// Miner is the URL of the miner that submitted the transaction.
	Miner string `json:"miner,omitempty"`
	// BlockHash is the hash of the block where the transaction was included.
	BlockHash string `json:"blockHash,omitempty"`
	// BlockHeight is the height of the block where the transaction was included.
	BlockHeight int64 `json:"blockHeight,omitempty"`
	// ExtraInfo is the extra info returned by the broadcaster.
	ExtraInfo string `json:"extraInfo,omitempty"`
	// Status is the status of the response.
	Status int `json:"status,omitempty"`
	// Title is the title of the response.
	Title string `json:"title,omitempty"`
	// TxStatus is the status of the transaction.
	TxStatus TxStatus `json:"txStatus,omitempty"`
}

// SubmitTxResponse is the response returned by the SubmitBatchTransactions method.
type SubmitBatchTxResponse struct {
	// Miner is the URL of the miner that submitted the batch transactions.
	Miner string `json:"miner"`
	// SubmitTxResponses are the responses returned by the miner who submitted the batch transactions.
	SubmitTxResponses []*SubmitTxResponse
}
