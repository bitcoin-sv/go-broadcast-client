package broadcast

import (
	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-api"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/composite"
)

type ClientBuilder struct {
	factories []composite.BroadcastFactory
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{}
}

func (cb *ClientBuilder) WithArc(config ArcClientConfig) *ClientBuilder {
	cb.factories = append(cb.factories, func() broadcast_api.Client {
		return arc.NewArcClient(&config, nil)
	})
	return cb
}

func (cb *ClientBuilder) Build() broadcast_api.Client {
	if len(cb.factories) == 1 {
		return cb.factories[0]()
	}
	return composite.NewBroadcasterWithDefaultStrategy(cb.factories...)
}
