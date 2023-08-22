package broadcast

// SubmittedTx is the response returned by the miner from submitting transaction(s).
type SubmittedTx struct {
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

// SubmitTxResponse is the response returned by the SubmitTransaction method.
type SubmitTxResponse struct {
	BaseResponse
	*SubmittedTx
}

// SubmitTxResponse is the response returned by the SubmitBatchTransactions method.
type SubmitBatchTxResponse struct {
	BaseResponse
	// Transactions is a slice of responses returned from submitting the batch transactions.
	Transactions []*SubmittedTx `json:"transactions"`
}
