package acceptancetests

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/rs/zerolog"

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

var successfulSubmitBatchResponse = `
[
	{
    	"txStatus": "MINED"
	},
	{
    	"txStatus": "MINED"
	}
]
`

const mockHex = "transaction-data"

func submitUrl(base string) string {
	return requestUrl(base, "/v1/tx")
}

func TestSubmitTransaction(t *testing.T) {
	testLogger := zerolog.Nop()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	t.Run("Should successfully submit transaction using first of two ArcClients", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(2)

		httpmock.RegisterResponder("POST", submitUrl(urls[0]),
			httpmock.NewStringResponder(http.StatusOK, successfulSubmitResponse),
		)
		// first miner responded successfully, next one should be skipped
		httpmock.RegisterResponder("GET", submitUrl(urls[1]),
			httpmock.NewStringResponder(http.StatusOK, successfulSubmitResponse),
		)

		// when
		result, err := broadcaster.SubmitTransaction(context.Background(), &broadcast.Transaction{Hex: "transaction-data"})

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("Should return error if both ArcClients return errors", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		broadcaster := broadcast_client.Builder().
			WithHttpClient(httpClientMock).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}, &testLogger).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc2-api-url", Token: "arc2-token"}, &testLogger).
			Build()

		httpResponse := &http.Response{}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, errors.New("http error")).Once()
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, errors.New("http error")).Once()

		// when
		result, err := broadcaster.SubmitTransaction(context.Background(), &broadcast.Transaction{Hex: "transaction-data"})

		// then
		httpClientMock.AssertExpectations(t)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, broadcast.ErrAllBroadcastersFailed.Error())
	})

	t.Run("Should return error if one ArcClient returns error and the other returns invalid response", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		broadcaster := broadcast_client.Builder().
			WithHttpClient(httpClientMock).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}, &testLogger).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc2-api-url", Token: "arc2-token"}, &testLogger).
			Build()

		httpResponse := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("invalid-response"))}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, nil).Once()
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, errors.New("http error")).Once()

		// when
		result, err := broadcaster.SubmitTransaction(context.Background(), &broadcast.Transaction{Hex: "transaction-data"})

		// then
		httpClientMock.AssertExpectations(t)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Should successfully submit transaction if one ArcClient returns missing status and the other is successful", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		broadcaster := broadcast_client.Builder().
			WithHttpClient(httpClientMock).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}, &testLogger).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc2-api-url", Token: "arc2-token"}, &testLogger).
			Build()

		httpResponse1 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"blockHash": "hash"}`))}
		httpResponse2 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(successfulSubmitResponse))}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse1, nil).Once()
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse2, nil).Once()

		// when
		result, err := broadcaster.SubmitTransaction(context.Background(), &broadcast.Transaction{Hex: "transaction-data"})

		// then
		httpClientMock.AssertExpectations(t)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestSubmitBatchTransactions(t *testing.T) {
	t.Run("Should successfully submit batch of transactions using first of few ArcClients", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		broadcaster := buildBroadcastClient(2, httpClientMock)

		httpResponse1 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(successfulSubmitBatchResponse))}
		httpResponse2 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(successfulSubmitBatchResponse))}

		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse1, nil).Once()
		// first miner responded successfully, next one should be skipped
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse2, nil).Times(0)

		batch := []*broadcast.Transaction{
			{Hex: "transaction-0-data"},
			{Hex: "transaction-1-data"},
		}

		// when
		result, err := broadcaster.SubmitBatchTransactions(context.Background(), batch)

		// then
		httpClientMock.AssertExpectations(t)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Should successfully submit batch of transactions if one ArcClient returns missing status and the other is successful", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		broadcaster := buildBroadcastClient(2, httpClientMock)

		// response without status field
		httpResponse1 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`[{"blockHash": "hash"}, {"txStatus": "MINED"}`))}
		// valid response
		httpResponse2 := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(successfulSubmitBatchResponse))}

		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse1, nil).Once()
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse2, nil).Once()

		batch := []*broadcast.Transaction{
			{Hex: "transaction-0-data"},
			{Hex: "transaction-1-data"},
		}

		// when
		result, err := broadcaster.SubmitBatchTransactions(context.Background(), batch)

		// then
		httpClientMock.AssertExpectations(t)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Should return error if every ArcClients return errors", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		broadcaster := buildBroadcastClient(2, httpClientMock)

		httpResponse := &http.Response{}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, errors.New("http error")).Once()
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, errors.New("http error")).Once()

		batch := []*broadcast.Transaction{
			{Hex: "transaction-0-data"},
			{Hex: "transaction-1-data"},
		}

		// when
		result, err := broadcaster.SubmitBatchTransactions(context.Background(), batch)

		// then
		httpClientMock.AssertExpectations(t)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, broadcast.ErrAllBroadcastersFailed.Error())
	})

	t.Run("Should return error if one ArcClient returns error and the other returns invalid response", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		broadcaster := buildBroadcastClient(2, httpClientMock)

		httpResponse := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("invalid-response"))}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, nil).Once()
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).Return(httpResponse, errors.New("http error")).Once()

		batch := []*broadcast.Transaction{
			{Hex: "transaction-0-data"},
			{Hex: "transaction-1-data"},
		}

		// when
		result, err := broadcaster.SubmitBatchTransactions(context.Background(), batch)

		// then
		httpClientMock.AssertExpectations(t)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func buildBroadcastClient(acrClients uint, httpClient *arc.MockHttpClient) broadcast.Client {
	builder := broadcast_client.Builder().
		WithHttpClient(httpClient)
	testLogger := zerolog.Nop()

	for i := uint(0); i < acrClients; i++ {
		arc := broadcast_client.ArcClientConfig{
			APIUrl: fmt.Sprintf("http://arc%d-api-url", i),
			Token:  fmt.Sprintf("arc%d-token", i),
		}
		builder = builder.WithArc(arc, &testLogger)
	}

	return builder.Build()
}
