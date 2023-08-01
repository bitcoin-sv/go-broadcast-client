package acceptancetests

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
	arc3 "github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/composite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var firstSuccessfulResponse = `
{
    "blockHash": "hash",
    "txStatus": "MINED"
}
`

var secondSuccessfulResponse = `
{
    "blockHash": "hash",
    "txStatus": "CONFIRMED"
}
`

func TestQueryTransaction(t *testing.T) {
	t.Run("Should successfully query from multiple ArcClients", func(t *testing.T) {
		// given
		httpClientMock := &arc3.MockHttpClient{}
		arc1 := arc3.NewArcClient(&broadcast_client.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}, httpClientMock)
		arc2 := arc3.NewArcClient(&broadcast_client.ArcClientConfig{APIUrl: "http://arc2-api-url", Token: "arc2-token"}, httpClientMock)
		broadcaster := composite.NewBroadcaster(composite.DefaultStrategy, arc3.CreateMockArcFactory(arc1), arc3.CreateMockArcFactory(arc2))
		httpResponse1 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(firstSuccessfulResponse))}
		httpResponse2 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(secondSuccessfulResponse))}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse1, nil).Once()
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse2, nil).Once()

		// when
		result, err := broadcaster.QueryTransaction(context.Background(), "txID")

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Should return error if all ArcClients return errors", func(t *testing.T) {
		// given
		httpClientMock := &arc3.MockHttpClient{}
		arc1 := arc3.NewArcClient(&broadcast_client.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}, httpClientMock)
		arc2 := arc3.NewArcClient(&broadcast_client.ArcClientConfig{APIUrl: "http://arc2-api-url", Token: "arc2-token"}, httpClientMock)
		broadcaster := composite.NewBroadcaster(composite.DefaultStrategy, arc3.CreateMockArcFactory(arc1), arc3.CreateMockArcFactory(arc2))
		httpResponse := &http.Response{}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, errors.New("http error")).Twice()

		// when
		result, err := broadcaster.QueryTransaction(context.Background(), "txID")

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, broadcast.ErrAllBroadcastersFailed.Error())
	})

	t.Run("Should successfully query from single ArcClient", func(t *testing.T) {
		// given
		httpClientMock := &arc3.MockHttpClient{}
		arc1 := arc3.NewArcClient(&broadcast_client.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}, httpClientMock)
		broadcaster := composite.NewBroadcaster(composite.DefaultStrategy, arc3.CreateMockArcFactory(arc1))
		httpResponse1 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(firstSuccessfulResponse))}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse1, nil).Once()

		// when
		result, err := broadcaster.QueryTransaction(context.Background(), "txID")

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Should return error if single ArcClient returns error", func(t *testing.T) {
		// given
		httpClientMock := &arc3.MockHttpClient{}
		arc1 := arc3.NewArcClient(&broadcast_client.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}, httpClientMock)
		broadcaster := composite.NewBroadcaster(composite.DefaultStrategy, arc3.CreateMockArcFactory(arc1))
		httpResponse := &http.Response{}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, errors.New("http error")).Once()

		// when
		result, err := broadcaster.QueryTransaction(context.Background(), "txID")

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
