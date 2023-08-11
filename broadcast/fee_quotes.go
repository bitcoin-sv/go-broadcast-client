package broadcast

type PolicyQuoteResponse struct {
	Miner     string         `json:"miner"`
	Policy    PolicyResponse `json:"policy"`
	Timestamp string         `json:"timestamp"`
}

type FeeQuote struct {
	Miner     string            `json:"miner"`
	MiningFee MiningFeeResponse `json:"miningFee"`
	Timestamp string            `json:"timestamp"`
}

type PolicyResponse struct {
	MaxScriptSizePolicy    int64             `json:"maxscriptsizepolicy"`
	MaxTxSigOpsCountPolicy int64             `json:"maxtxsigopscountspolicy"`
	MaxTxSizePolicy        int64             `json:"maxtxsizepolicy"`
	MiningFee              MiningFeeResponse `json:"miningFee"`
}

type MiningFeeResponse struct {
	Bytes    int64 `json:"bytes"`
	Satoshis int64 `json:"satoshis"`
}
