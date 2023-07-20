package factory

import (
	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/arc"
	"github.com/bitcoin-sv/go-broadcast-client/config"
)

func NewBroadcastClient(factories ...broadcast.BroadcastFactory) broadcast.Broadcaster {
	return broadcast.NewCompositeBroadcaster(config.DefaultStrategy, factories...)
}

func WithArc(config config.ArcClientConfig) broadcast.BroadcastFactory {
	return func() broadcast.Broadcaster {
		return arc.NewArcClient(config)
	}
}

func WithMultiminers(factories ...broadcast.BroadcastFactory) broadcast.BroadcastFactory {
	return func() broadcast.Broadcaster {
		return broadcast.NewCompositeBroadcaster(config.DefaultStrategy, factories...)
	}
}
