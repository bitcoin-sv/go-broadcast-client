package broadcast

import "github.com/bitcoin-sv/go-broadcast-client/internal/arc"

type ClientBuilder struct {
	factories []BroadcastFactory
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{}
}

func (cb *ClientBuilder) WithArc(config ArcClientConfig) *ClientBuilder {
	cb.factories = append(cb.factories, func() Broadcaster {
		return arc.NewArcClient(config, nil)
	})
	return cb
}

func (cb *ClientBuilder) Build() Broadcaster {
	if len(cb.factories) == 1 {
		return cb.factories[0]()
	}
	return NewCompositeBroadcaster(DefaultStrategy, cb.factories...)
}
