package broadcast_client

import (
	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/composite"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/httpclient"
)

type builder struct {
	factories []composite.BroadcastFactory
	client    httpclient.HTTPInterface
}

// Builder is used to prepare the broadcast client. It is recommended to use that builder for creating the broadcast client.
func Builder() *builder {
	return &builder{}
}

// WithHttpClient sets the http client to be used by the broadcast client. It requires a httpclient.HTTPInterface to be passed.
func (cb *builder) WithHttpClient(client httpclient.HTTPInterface) *builder {
	cb.client = client
	return cb
}

// WithArc sets up the connection of the broadcast client to the Arc service using the provided ArcClientConfig.
// This method can be called multiple times with different ArcClientConfigurations to establish connections to multiple Arc instances.
func (cb *builder) WithArc(config ArcClientConfig) *builder {
	cb.factories = append(cb.factories, func() broadcast_api.Client {
		return arc.NewArcClient(&config, cb.client)
	})
	return cb
}

// Build builds the broadcast client based on the provided configuration.
func (cb *builder) Build() broadcast_api.Client {
	if len(cb.factories) == 1 {
		return cb.factories[0]()
	}
	return composite.NewBroadcasterWithDefaultStrategy(cb.factories...)
}
