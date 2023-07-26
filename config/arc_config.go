package config

var ArcQueryTxRoute = "/v1/tx/"
var ArcPolicyQuoteRoute = "/v1/policy"
var SubmitTxRoute = "/v1/tx"
var SubmitTxsRoute = "/v1/txs"

type ArcClientConfig struct {
	APIUrl string
	Token  string
}
