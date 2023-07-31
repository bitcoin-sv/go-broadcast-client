package factory

import (
	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/config"
	"github.com/bitcoin-sv/go-broadcast-client/internal/arc"
)

func NewBroadcastClient(factories ...broadcast.BroadcastFactory) broadcast.Broadcaster {
	if len(factories) == 1 {
		return factories[0]()
	}

	return broadcast.NewCompositeBroadcaster(config.DefaultStrategy, factories...)
}

func WithArc(config config.ArcClientConfig) broadcast.BroadcastFactory {
	return func() broadcast.Broadcaster {
		return arc.NewArcClient(config)
	}
}
