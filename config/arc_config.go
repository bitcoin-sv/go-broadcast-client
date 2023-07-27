package config

var ArcQueryTxRoute = "/v1/tx/"
var ArcPolicyQuoteRoute = "/v1/policy"
var ArcSubmitTxRoute = "/v1/tx"
var ArcSubmitTxsRoute = "/v1/txs"

type ArcClientConfig struct {
	APIUrl string
	Token  string
}
