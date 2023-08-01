package broadcast_api

// Transaction is the body contents in the "submit transaction" request
type Transaction struct {
	CallBackEncryption string   `json:"callBackEncryption,omitempty"`
	CallBackToken      string   `json:"callBackToken,omitempty"`
	CallBackURL        string   `json:"callBackUrl,omitempty"`
	DsCheck            bool     `json:"dsCheck,omitempty"`
	MerkleFormat       string   `json:"merkleFormat,omitempty"`
	MerkleProof        bool     `json:"merkleProof,omitempty"`
	RawTx              string   `json:"rawtx"`
	WaitForStatus      TxStatus `json:"waitForStatus,omitempty"`
}
