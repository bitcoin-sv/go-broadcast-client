package arc

import (
	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/config"
	"github.com/bitcoin-sv/go-broadcast-client/internal/httpclient"
)

type ArcClient struct {
	apiURL     string
	token      string
	httpClient httpclient.HTTPInterface
}

func NewArcClient(config config.ArcClientConfig) broadcast.Broadcaster {
	httpClient := httpclient.NewHttpClient()

	return &ArcClient{
		apiURL:     config.APIUrl,
		token:      config.Token,
		httpClient: httpClient,
	}
}
