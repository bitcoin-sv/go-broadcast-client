package broadcast_client

import (
	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc/mocks"
)

// MockType is an enum that is used as parameter to WithMockArc
// client builder in order to create different types of mock.
type MockType int

const (
	MockSuccess MockType = iota
	MockFailure
	MockTimeout
)

// WithMockArc creates a mock client for testing purposes. It takes mock type as argument
// and creates a mock client that satisfies the client interface with methods that return
// success or specific error based on this mock type argument.
func (cb *builder) WithMockArc(config ArcClientConfig, mockType MockType) *builder {
	var clientToReturn broadcast_api.Client

	switch mockType {
	case MockSuccess:
		clientToReturn = mocks.NewArcClientMock()
		break
	case MockFailure:
		clientToReturn = mocks.NewArcClientMockFailure()
		break
	case MockTimeout:
		clientToReturn = mocks.NewArcClientMockTimeout()
		break
	default:
		clientToReturn = mocks.NewArcClientMock()
	}

	cb.factories = append(cb.factories, func() broadcast_api.Client {
		return clientToReturn
	})
	return cb
}
