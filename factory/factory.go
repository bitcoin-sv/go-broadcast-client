package factory

import (
	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/config"
	"github.com/bitcoin-sv/go-broadcast-client/internal/arc"
)

type ClientBuilder struct {
	factories []broadcast.BroadcastFactory
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{}
}

func NewBroadcastClient(factories ...broadcast.BroadcastFactory) broadcast.Broadcaster {
	if len(factories) == 1 {
		return factories[0]()
	}

	return broadcast.NewCompositeBroadcaster(config.DefaultStrategy, factories...)
}

func (cb *ClientBuilder) WithArc(config config.ArcClientConfig) *ClientBuilder {
	cb.factories = append(cb.factories, func() broadcast.Broadcaster {
		return arc.NewArcClient(config, nil)
	})
	return cb
}

func (cb *ClientBuilder) Build() broadcast.Broadcaster {
	if len(cb.factories) == 1 {
		return cb.factories[0]()
	}
	return broadcast.NewCompositeBroadcaster(config.DefaultStrategy, cb.factories...)
}
