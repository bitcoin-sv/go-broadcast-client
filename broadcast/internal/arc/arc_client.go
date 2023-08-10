// Package arc provides the implementation of the arc client.
package arc

import (
	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/httpclient"
)

const (
	arcQueryTxRoute        = "/v1/tx/"
	arcPolicyQuoteRoute    = "/v1/policy"
	arcSubmitTxRoute       = "/v1/tx"
	arcSubmitBatchTxsRoute = "/v1/txs"
)

// Config is the interface that wraps the basic functions for the arc client configuration.
type Config interface {
	// GetApiUrl returns the arc api url.
	GetApiUrl() string
	// GetToken returns the arc api token.
	GetToken() string
}

// ArcClient is the implementation of the arc client.
type ArcClient struct {
	apiURL string
	token  string
	// HTTPClient is the http client used to make requests to the arc api.
	HTTPClient httpclient.HTTPInterface
}

// NewArcClient returns a new instance of the arc client.
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
