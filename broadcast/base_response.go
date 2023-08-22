package broadcast

type BaseResponse struct {
	// Miner is the URL of the miner that returned the response.
	Miner string `json:"miner"`
}
