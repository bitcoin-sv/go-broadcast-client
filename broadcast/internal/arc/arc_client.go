package arc

import (
	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/httpclient"
	"github.com/rs/zerolog"
)

const (
	arcQueryTxRoute        = "/v1/tx/"
	arcPolicyQuoteRoute    = "/v1/policy"
	arcSubmitTxRoute       = "/v1/tx"
	arcSubmitBatchTxsRoute = "/v1/txs"
)

type Config interface {
	GetApiUrl() string
	GetToken() string
	GetDeploymentID() string
}

type ArcClient struct {
	apiURL       string
	token        string
	deploymentID string
	HTTPClient   httpclient.HTTPInterface
	Logger       *zerolog.Logger
}

func NewArcClient(config Config, client httpclient.HTTPInterface, log *zerolog.Logger) broadcast_api.Client {
	if client == nil {
		client = httpclient.NewHttpClient()
	}

	arcClient := &ArcClient{
		apiURL:       config.GetApiUrl(),
		token:        config.GetToken(),
		deploymentID: config.GetDeploymentID(),
		HTTPClient:   client,
		Logger:       log,
	}

	log.Debug().Msgf("Created new arc client with api url: %s", arcClient.apiURL)
	return arcClient
}
