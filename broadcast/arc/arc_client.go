package arc

import (
	"context"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/config"
)

type ArcClient struct {
	apiURL string
	token  string
}

func NewArcClient(config config.ArcClientConfig) broadcast.Broadcaster {
	return &ArcClient{apiURL: config.APIUrl, token: config.Token}
}

func (a *ArcClient) BestQuote(ctx context.Context, feeCategory, feeType string) error {
	// implementacja
	return nil
}
