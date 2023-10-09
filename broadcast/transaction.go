package broadcast

// Transaction is the body contents in the "submit transaction" request.
type Transaction struct {
	// Hex is the transaction hex string.
	Hex string `json:"hex"`
}

// BasicTxResponse is a base type for query and submit transaction response.
type BaseTxResponse struct {
	// BlockHash is the hash of the block where the transaction was included.
	BlockHash string `json:"blockHash,omitempty"`
	// BlockHeight is the height of the block where the transaction was included.
	BlockHeight int64 `json:"blockHeight,omitempty"`
	// ExtraInfo provides extra information for given transaction.
	ExtraInfo string `json:"extraInfo,omitempty"`
	// MerklePath is the Merkle path used to calculate Merkle root of the block in which the transaction was included.
	MerklePath string `json:"merklePath,omitempty"`
	// Timestamp is the timestamp of the block where the transaction was included.
	Timestamp string `json:"timestamp,omitempty"`
	// TxStatus is the status of the transaction.
	TxStatus TxStatus `json:"txStatus,omitempty"`
	// TxID is the transaction id.
	TxID string `json:"txid,omitempty"`
}

// QueryTxResponse is the response returned by the QueryTransaction method.
type QueryTxResponse struct {
	BaseResponse
	BaseTxResponse
}

// SubmittedTx is the internal response returned by the miner from submitting transaction(s).
type SubmittedTx struct {
	BaseTxResponse
	// Status is the status of the response.
	Status int `json:"status,omitempty"`
	// Title is the title of the response.
	Title string `json:"title,omitempty"`
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
