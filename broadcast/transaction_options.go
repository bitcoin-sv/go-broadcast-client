package broadcast

// TransactionOptFunc defines an optional arguments that can be passed to the SubmitTransaction method.
type TransactionOptFunc func(o *TransactionOpts)

// TransactionFormat is the format of transaction being submitted.
type TransactionFormat int

const (
	EfFormat TransactionFormat = iota
	BeefFormat
	RawTxFormat
)

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

// WithCallback allow you to get the callback from the node when the transaction is mined,
// and receive the transaction details and status.
func WithCallback(callbackURL string, callbackToken ...string) TransactionOptFunc {
	return func(o *TransactionOpts) {
		o.CallbackURL = callbackURL
		if len(callbackToken) > 0 {
			o.CallbackToken = callbackToken[0]
		}
	}
}

// WithMerkleProofFromCallback it's an option that indicates if the merkle proof should be returned in callback.
func WithMerkleProofFromCallback() TransactionOptFunc {
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
