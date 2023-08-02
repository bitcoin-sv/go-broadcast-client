package acceptancetests

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	broadcast_client "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var successfulSubmitResponse = `
{
    "txStatus": "MINED"
}
`

func TestSubmitTransaction(t *testing.T) {
	t.Run("Should successfully submit transaction using two ArcClients", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		broadcaster := broadcast_client.Builder().
			WithHttpClient(httpClientMock).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc2-api-url", Token: "arc2-token"}).
			Build()

		httpResponse1 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(successfulSubmitResponse))}
		httpResponse2 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(successfulSubmitResponse))}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse1, nil).Once()
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse2, nil).Once()

		// when
		result, err := broadcaster.SubmitTransaction(context.Background(), &broadcast.Transaction{RawTx: "transaction-data"})

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Should return error if both ArcClients return errors", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		broadcaster := broadcast_client.Builder().
			WithHttpClient(httpClientMock).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc2-api-url", Token: "arc2-token"}).
			Build()

		httpResponse := &http.Response{}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, errors.New("http error")).Once()
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, errors.New("http error")).Once()

		// when
		result, err := broadcaster.SubmitTransaction(context.Background(), &broadcast.Transaction{RawTx: "transaction-data"})

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, broadcast.ErrAllBroadcastersFailed.Error())
	})

	t.Run("Should return error if one ArcClient returns error and the other returns invalid response", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		broadcaster := broadcast_client.Builder().
			WithHttpClient(httpClientMock).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc2-api-url", Token: "arc2-token"}).
			Build()

		httpResponse := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("invalid-response"))}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, nil).Once()
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, errors.New("http error")).Once()

		// when
		result, err := broadcaster.SubmitTransaction(context.Background(), &broadcast.Transaction{RawTx: "transaction-data"})

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should return error if one ArcClient returns missing status and the other is successful", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		broadcaster := broadcast_client.Builder().
			WithHttpClient(httpClientMock).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc2-api-url", Token: "arc2-token"}).
			Build()

		httpResponse1 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(successfulSubmitResponse))}
		httpResponse2 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"blockHash": "hash"}`))}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse1, nil).Once()
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse2, nil).Once()

		// when
		result, err := broadcaster.SubmitTransaction(context.Background(), &broadcast.Transaction{RawTx: "transaction-data"})

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
