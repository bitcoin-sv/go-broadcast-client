package arc

import (
	"context"
	"net/http"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/httpclient"
	"github.com/stretchr/testify/mock"
)

// MockHttpClient is a mock implementation of the http client.
type MockHttpClient struct {
	mock.Mock
}

// DoRequest is the mock implementation of the http client DoRequest function required by HTTPInterface.
func (m *MockHttpClient) DoRequest(ctx context.Context, pld httpclient.HTTPRequest) (*http.Response, error) {
	args := m.Called(ctx, pld)
	return args.Get(0).(*http.Response), args.Error(1)
}
