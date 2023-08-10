package broadcast

// Transaction is the body contents in the "submit transaction" request
type Transaction struct {
	// CallbackEncryption is the encryption method used for the callback (to set in the callback header).
	CallBackEncryption string `json:"callBackEncryption,omitempty"`
	// CallbackToken is the token used for the callback (to set in the callback header).
	CallBackToken string `json:"callBackToken,omitempty"`
	// CallbackURL is the URL used for the callback (to set in the callback header).
	CallBackURL string `json:"callBackUrl,omitempty"`
	// DsCheck is the double spend check flag.
	DsCheck bool `json:"dsCheck,omitempty"`
	// MerkleFormat is the requested merkle format (to set in the merkle format header).
	MerkleFormat string `json:"merkleFormat,omitempty"`
	// MerkleProof is the merkle proof flag - if merkle proof should be included in the response (to set in the merkle proof header).
	MerkleProof bool `json:"merkleProof,omitempty"`
	// RawTx is the raw transaction string -> required.
	RawTx string `json:"rawtx"`
	// WaitForStatus is the status to wait for with the callback (to set in the wait for status header).
	WaitForStatus TxStatus `json:"waitForStatus,omitempty"`
}
