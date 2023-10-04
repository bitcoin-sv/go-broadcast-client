package broadcast

// Transaction is the body contents in the "submit transaction" request
type Transaction struct {
	// TxHex is the transaction hex string.
	Hex string `json:"hex"`
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
	// BeefFormat is a flag which indicates that the transaction given for submitting is in BEEF Format.
	BeefFormat bool
	// RawFormat is a flag which indicates that the transaction given for submitting is in Raw Tx Format.
	RawFormat bool
}

// WithCallback allow you to get the callback from the node when the transaction is mined, 
// and receive the transaction details and status.
func WithCallback(callbackURL string, callbackToken ...string) TransactionOptFunc {
	return func(o *TransactionOpts) {
		o.CallbackToken = callbackURL
		if len(callbackToken) > 0 {
			o.CallbackToken = callbackToken[0]
		}
	}
}

// WithMerkleProof will return merkle proof from Arc.
func WithMerkleProof() TransactionOptFunc {
	return func(o *TransactionOpts) {
		o.MerkleProof = true
	}
}

// WithWaitForStatus will allow you to return the result only
// when the transaction reaches the status you set.
func WithWaitForStatus(status TxStatus) TransactionOptFunc {
	return func(o *TransactionOpts) {
		o.WaitForStatus = status
	}
}

// WithBeefFormat will accept your transaction in BEEF format 
// and decode it for a proper format acceptable by Arc.
func WithBeefFormat() TransactionOptFunc {
	return func(o *TransactionOpts) {
		o.BeefFormat = true
		o.RawFormat = false
	}
}

// WithRawFormat will accept your transaction in RawTx format 
// and encode it for a proper format acceptable by Arc.
// Deprecated: This function will be depreacted soon.
// Only EF and BEEF format will be acceptable.
func WithRawFormat() TransactionOptFunc {
	return func(o *TransactionOpts) {
		o.BeefFormat = false
		o.RawFormat = true
	}
}
