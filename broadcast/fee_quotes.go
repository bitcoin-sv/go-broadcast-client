package broadcast

// PolicyQuoteResponse is the response returned by the GetPolicyQuote method.
type PolicyQuoteResponse struct {
	BaseResponse
	// Policy is a detailed policy of the miner.
	Policy PolicyResponse `json:"policy"`
	// Timestamp is the timestamp of the policy quote response received.
	Timestamp string `json:"timestamp"`
}

// FeeQuote is the response returned by the GetFeeQuote method.
type FeeQuote struct {
	BaseResponse
	// MiningFee is expressed by number of satoshis per number of bytes.
	MiningFee MiningFeeResponse `json:"miningFee"`
	// Timestamp is the timestamp of the fee quote response received.
	Timestamp string `json:"timestamp"`
}

type PolicyResponse struct {
	// MaxScriptSizePolicy is the maximum number of bytes of the script in a single transaction.
	MaxScriptSizePolicy int64 `json:"maxscriptsizepolicy"`
	// MaxTxSigOpsCountPolicy is the maximum allowed number of signature operations in a single transaction
	MaxTxSigOpsCountPolicy int64 `json:"maxtxsigopscountspolicy"`
	// MaxTxSizePolicy is the maximum size in bytes of a single transaction.
	MaxTxSizePolicy int64 `json:"maxtxsizepolicy"`
	// MiningFee is expressed by number of satoshis per number of bytes.
	MiningFee MiningFeeResponse `json:"miningFee"`
}

type MiningFeeResponse struct {
	// Bytes is the number of bytes per which a number of satoshis is charged as transaction fee.
	Bytes int64 `json:"bytes"`
	// Satoshis is the number of satoshis charged as transaction fee per the number of bytes.
	Satoshis int64 `json:"satoshis"`
}
