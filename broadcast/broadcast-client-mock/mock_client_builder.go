package broadcast_client_mock

import (
	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc/mocks"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/composite"
)

// MockType is an enum that is used as parameter to WithMockArc
// client builder in order to create different types of mock.
type MockType int

const (
	MockSuccess MockType = iota
	MockFailure
	MockTimeout
	MockNilQueryTxResp
)

type builder struct {
	factories []composite.BroadcastFactory
}

// Builder is used to prepare the mock broadcast client. It is recommended
// to use that builder for creating the mock broadcast client.
func Builder() *builder {
	return &builder{}
}

// WithMockArc creates a mock client for testing purposes. It takes mock type as argument
// and creates a mock client that satisfies the client interface with methods that return
// success or specific error based on this mock type argument. It allows for creating
// multiple mock clients.
func (cb *builder) WithMockArc(mockType MockType) *builder {
	var clientToReturn broadcast_api.Client

	switch mockType {
	case MockSuccess:
		clientToReturn = mocks.NewArcClientMock()
	case MockFailure:
		clientToReturn = mocks.NewArcClientMockFailure()
	case MockTimeout:
		clientToReturn = mocks.NewArcClientMockTimeout()
	case MockNilQueryTxResp:
		clientToReturn = mocks.NewArcClientMockNilQueryTxResp()
	default:
		clientToReturn = mocks.NewArcClientMock()
	}

	cb.factories = append(cb.factories, func() broadcast_api.Client {
		return clientToReturn
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
