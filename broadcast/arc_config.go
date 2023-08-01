package broadcast

const (
	ArcQueryTxRoute     = "/v1/tx/"
	ArcPolicyQuoteRoute = "/v1/policy"
	ArcSubmitTxRoute    = "/v1/tx"
	ArcSubmitTxsRoute   = "/v1/txs"
)

type ArcClientConfig struct {
	APIUrl string
	Token  string
}
