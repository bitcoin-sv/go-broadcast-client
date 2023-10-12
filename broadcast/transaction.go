package broadcast

// Transaction is the body contents in the "submit transaction" request
type Transaction struct {
	// RawTx is the raw transaction hex string.
	RawTx string `json:"rawtx"`
}

// TransactionOptFunc defines an optional arguments that can be passed to the SubmitTransaction method.
type TransactionOptFunc func(o *TransactionOpts)

// TransactionOpts is a struct that holds optional arguments that can be passed to the SubmitTransaction method.
type TransactionOpts struct {
	// CallbackURL is the URL that will be called when the transaction status changes.
	CallbackURL string
	// CallbackToken is the token that will be sent in the callback request.
	CallbackToken string
	// MerkleProof is a flag that indicates if the merkle proof should be returned in the submit transaction response.
	MerkleProof bool
	// WaitForStatus is the status that the callback request will wait for.
	WaitForStatus TxStatus
}

func WithCallback(callbackURL string, callbackToken ...string) TransactionOptFunc {
	return func(o *TransactionOpts) {
		o.CallbackToken = callbackURL
		if len(callbackToken) > 0 {
			o.CallbackToken = callbackToken[0]
		}
	}
}

func WithMerkleProof() TransactionOptFunc {
	return func(o *TransactionOpts) {
		o.MerkleProof = true
	}
}

func WithWaitForStatus(status TxStatus) TransactionOptFunc {
	return func(o *TransactionOpts) {
		o.WaitForStatus = status
	}
}
