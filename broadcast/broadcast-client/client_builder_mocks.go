package broadcast_client

import (
	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc/mocks"
)

type MockType int

const (
	MockSuccess MockType = iota + 1
	MockFailure
	MockTimeout
	MockAlreadyOnChain
	MockInMempool
)

func (cb *builder) WithMockArc(config ArcClientConfig, mockType MockType) *builder {
	var clientToReturn broadcast_api.Client

	switch mockType {
	case MockSuccess:
		clientToReturn = mocks.NewArcClientMock()
		break
	case MockFailure:
		clientToReturn = mocks.NewArcClientMockFailure()
		break
	case MockAlreadyOnChain:
		clientToReturn = mocks.NewArcClientMockOnChain()
		break
	default:
		clientToReturn = mocks.NewArcClientMock()
	}

	cb.factories = append(cb.factories, func() broadcast_api.Client {
		return clientToReturn
	})
	return cb
}
