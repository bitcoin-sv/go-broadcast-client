package broadcast_client

import (
	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/composite"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/httpclient"
)

type builder struct {
	factories []composite.BroadcastFactory
}

func Builder() *builder {
	return &builder{}
}

func (cb *builder) WithArc(config ArcClientConfig, client httpclient.HTTPInterface) *builder {
	cb.factories = append(cb.factories, func() broadcast_api.Client {
		return arc.NewArcClient(&config, client)
	})
	return cb
}

func (cb *builder) Build() broadcast_api.Client {
	if len(cb.factories) == 1 {
		return cb.factories[0]()
	}
	return composite.NewBroadcasterWithDefaultStrategy(cb.factories...)
}
