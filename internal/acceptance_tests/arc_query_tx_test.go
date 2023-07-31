package acceptancetests

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/config"
	"github.com/bitcoin-sv/go-broadcast-client/internal/arc"
	"github.com/bitcoin-sv/go-broadcast-client/shared"
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
		httpClientMock := &arc.MockHttpClient{}
		arc1 := arc.NewArcClient(config.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}, httpClientMock)
		arc2 := arc.NewArcClient(config.ArcClientConfig{APIUrl: "http://arc2-api-url", Token: "arc2-token"}, httpClientMock)
		broadcaster := broadcast.NewCompositeBroadcaster(config.DefaultStrategy, arc.CreateMockArcFactory(arc1), arc.CreateMockArcFactory(arc2))
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
		httpClientMock := &arc.MockHttpClient{}
		arc1 := arc.NewArcClient(config.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}, httpClientMock)
		arc2 := arc.NewArcClient(config.ArcClientConfig{APIUrl: "http://arc2-api-url", Token: "arc2-token"}, httpClientMock)
		broadcaster := broadcast.NewCompositeBroadcaster(config.DefaultStrategy, arc.CreateMockArcFactory(arc1), arc.CreateMockArcFactory(arc2))
		httpResponse := &http.Response{}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, errors.New("http error")).Twice()

		// when
		result, err := broadcaster.QueryTransaction(context.Background(), "txID")

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, shared.ErrAllBroadcastersFailed.Error())
	})

	t.Run("Should successfully query from single ArcClient", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		arc1 := arc.NewArcClient(config.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}, httpClientMock)
		broadcaster := broadcast.NewCompositeBroadcaster(config.DefaultStrategy, arc.CreateMockArcFactory(arc1))
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
		httpClientMock := &arc.MockHttpClient{}
		arc1 := arc.NewArcClient(config.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}, httpClientMock)
		broadcaster := broadcast.NewCompositeBroadcaster(config.DefaultStrategy, arc.CreateMockArcFactory(arc1))
		httpResponse := &http.Response{}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, errors.New("http error")).Once()

		// when
		result, err := broadcaster.QueryTransaction(context.Background(), "txID")

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
