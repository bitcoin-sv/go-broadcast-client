package arc

import (
	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/httpclient"
)

const (
	arcQueryTxRoute     = "/v1/tx/"
	arcPolicyQuoteRoute = "/v1/policy"
	arcSubmitTxRoute    = "/v1/tx"
	arcSubmitTxsRoute   = "/v1/txs"
)

type Config interface {
	GetApiUrl() string
	GetToken() string
}

type ArcClient struct {
	apiURL     string
	token      string
	HTTPClient httpclient.HTTPInterface
}

func NewArcClient(config Config, client httpclient.HTTPInterface) broadcast_api.Client {
	if client == nil {
		client = httpclient.NewHttpClient()
	}
	return &ArcClient{
		apiURL:     config.GetApiUrl(),
		token:      config.GetToken(),
		HTTPClient: client,
	}
}
