package arc

import (
	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/internal/httpclient"
)

type ArcClient struct {
	apiURL     string
	token      string
	HTTPClient httpclient.HTTPInterface
}

func NewArcClient(config broadcast.ArcClientConfig, client httpclient.HTTPInterface) broadcast.Broadcaster {
	if client == nil {
		client = httpclient.NewHttpClient()
	}
	return &ArcClient{
		apiURL:     config.APIUrl,
		token:      config.Token,
		HTTPClient: client,
	}
}
