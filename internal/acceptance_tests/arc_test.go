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
    "txStatus": "MINED",
}
`

var secondSuccessfulResponse = `
{
    "blockHash": "hash",
    "txStatus": "CONFIRMED",
}
`

func TestQueryTransaction(t *testing.T) {
	t.Run("Happy Path - Query success from arc1 and arc2", func(t *testing.T) {
		httpClientMock := &arc.MockHttpClient{}

		arc1 := arc.NewArcClient(config.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"})
		arc2 := arc.NewArcClient(config.ArcClientConfig{APIUrl: "http://arc2-api-url", Token: "arc2-token"})

		broadcaster := broadcast.NewCompositeBroadcaster(config.DefaultStrategy, createMockArcFactory(arc1.(*arc.ArcClient)), createMockArcFactory(arc2.(*arc.ArcClient)))

		arc1.(*arc.ArcClient).HTTPClient = httpClientMock
		arc2.(*arc.ArcClient).HTTPClient = httpClientMock

		httpResponse1 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(firstSuccessfulResponse))}
		httpResponse2 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(secondSuccessfulResponse))}

		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse1, nil).Once()
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse2, nil).Once()

		result, err := broadcaster.QueryTransaction(context.Background(), "txID")
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Sad Path - All ArcClients return errors", func(t *testing.T) {
		httpClientMock := &arc.MockHttpClient{}

		arc1 := arc.NewArcClient(config.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"})
		arc2 := arc.NewArcClient(config.ArcClientConfig{APIUrl: "http://arc2-api-url", Token: "arc2-token"})

		broadcaster := broadcast.NewCompositeBroadcaster(config.DefaultStrategy, createMockArcFactory(arc1.(*arc.ArcClient)), createMockArcFactory(arc2.(*arc.ArcClient)))

		arc1.(*arc.ArcClient).HTTPClient = httpClientMock
		arc2.(*arc.ArcClient).HTTPClient = httpClientMock

		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(nil, errors.New("http error")).Twice()

		result, err := broadcaster.QueryTransaction(context.Background(), "txID")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, shared.ErrAllBroadcastersFailed.Error())
	})

	t.Run("Happy Path - Query success from single arc", func(t *testing.T) {
		httpClientMock := &arc.MockHttpClient{}

		arc1 := arc.NewArcClient(config.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"})

		broadcaster := broadcast.NewCompositeBroadcaster(config.DefaultStrategy, createMockArcFactory(arc1.(*arc.ArcClient)))

		arc1.(*arc.ArcClient).HTTPClient = httpClientMock

		httpResponse1 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(firstSuccessfulResponse))}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse1, nil).Once()

		result, err := broadcaster.QueryTransaction(context.Background(), "txID")
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Sad Path - Query error from single arc", func(t *testing.T) {
		httpClientMock := &arc.MockHttpClient{}

		arc1 := arc.NewArcClient(config.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"})

		broadcaster := broadcast.NewCompositeBroadcaster(config.DefaultStrategy, createMockArcFactory(arc1.(*arc.ArcClient)))

		arc1.(*arc.ArcClient).HTTPClient = httpClientMock

		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(nil, errors.New("http error")).Once()

		result, err := broadcaster.QueryTransaction(context.Background(), "txID")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

}

func createMockArcFactory(arcClient *arc.ArcClient) broadcast.BroadcastFactory {
	return func() broadcast.Broadcaster {
		return arcClient
	}
}
