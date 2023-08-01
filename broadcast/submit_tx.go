package broadcast

type SubmitTxResponse struct {
	BlockHash   string   `json:"blockHash,omitempty"`
	BlockHeight int64    `json:"blockHeight,omitempty"`
	ExtraInfo   string   `json:"extraInfo,omitempty"`
	Status      int      `json:"status,omitempty"`
	Title       string   `json:"title,omitempty"`
	TxStatus    TxStatus `json:"txStatus,omitempty"`
}
