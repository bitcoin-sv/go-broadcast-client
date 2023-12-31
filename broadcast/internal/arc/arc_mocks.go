package arc

import (
	"context"
	"net/http"

	"github.com/bitcoin-sv/go-broadcast-client/httpclient"
	"github.com/stretchr/testify/mock"
)

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) DoRequest(ctx context.Context, pld httpclient.HTTPRequest) (*http.Response, error) {
	args := m.Called(ctx, pld)
	return args.Get(0).(*http.Response), args.Error(1)
}
