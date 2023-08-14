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
	broadcast_client "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc"
)

var firstSuccessfulPolicyResponse = `
{
	"policy": {
		"maxscriptsizepolicy": 100000000,
		"maxtxsigopscountspolicy": 4294967295,
		"maxtxsizepolicy": 100000000,
		"miningFee": {
			"bytes": 1000,
			"satoshis": 1
		}
	},
	"timestamp": "2023-08-10T13:49:07.308687569Z"
}
`

var secondSuccessfulPolicyResponse = `
{
	"policy": {
		"maxscriptsizepolicy": 100000000,
		"maxtxsigopscountspolicy": 4294967295,
		"maxtxsizepolicy": 100000000,
		"miningFee": {
			"bytes": 1000,
			"satoshis": 2
		}
	},
	"timestamp": "2023-08-10T13:49:07.308687569Z"
}
`

func TestPolicyQuote(t *testing.T) {
	t.Run("Should successfully query from multiple ArcClients", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		broadcaster := broadcast_client.Builder().
			WithHttpClient(httpClientMock).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc2-api-url", Token: "arc2-token"}).
			Build()

		httpResponse1 := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(firstSuccessfulPolicyResponse)),
		}
		httpResponse2 := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(secondSuccessfulPolicyResponse)),
		}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).
			Return(httpResponse1, nil).
			Once()
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).
			Return(httpResponse2, nil).
			Once()

		// when
		result, err := broadcaster.GetPolicyQuote(context.Background())

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Should return error if all ArcClients return errors", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		broadcaster := broadcast_client.Builder().
			WithHttpClient(httpClientMock).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc2-api-url", Token: "arc2-token"}).
			Build()

		httpResponse := &http.Response{}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).
			Return(httpResponse, errors.New("http error")).
			Twice()

		// when
		result, err := broadcaster.GetPolicyQuote(context.Background())

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, broadcast.ErrNoMinerResponse.Error())
	})

	t.Run("Should successfully query from single ArcClient", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		broadcaster := broadcast_client.Builder().
			WithHttpClient(httpClientMock).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}).
			Build()

		httpResponse1 := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(firstSuccessfulPolicyResponse)),
		}
		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).
			Return(httpResponse1, nil).
			Once()

		// when
		result, err := broadcaster.GetPolicyQuote(context.Background())

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Should return error if single ArcClient returns error", func(t *testing.T) {
		// given
		httpClientMock := &arc.MockHttpClient{}
		httpResponse := &http.Response{}
		broadcaster := broadcast_client.Builder().
			WithHttpClient(httpClientMock).
			WithArc(broadcast_client.ArcClientConfig{APIUrl: "http://arc1-api-url", Token: "arc1-token"}).
			Build()

		httpClientMock.On("DoRequest", mock.Anything, mock.Anything).
			Return(httpResponse, errors.New("http error")).
			Once()

		// when
		result, err := broadcaster.GetPolicyQuote(context.Background())

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
