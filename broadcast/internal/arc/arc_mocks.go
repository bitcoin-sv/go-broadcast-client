package arc

import (
	"context"
	"net/http"

	broadcast_api "github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/composite"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/httpclient"
	"github.com/stretchr/testify/mock"
)

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) DoRequest(ctx context.Context, pld httpclient.HTTPRequest) (*http.Response, error) {
	args := m.Called(ctx, pld)
	return args.Get(0).(*http.Response), args.Error(1)
}

func CreateMockArcFactory(arcClient broadcast_api.Client) composite.BroadcastFactory {
	return func() broadcast_api.Client {
		return arcClient
	}
}
