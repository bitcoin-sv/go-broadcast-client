package acceptancetests

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jarcoal/httpmock"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

const firstSuccessfulFeeQuoteResponse = `
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

const secondSuccessfulFeeQuoteResponse = `
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

func policyUrl(base string) string {
	return requestUrl(base, "/v1/policy")
}

func TestFeeQuote(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Should successfully query from multiple ArcClients", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(2)

		httpmock.RegisterResponder("GET", policyUrl(urls[0]),
			httpmock.NewStringResponder(http.StatusOK, firstSuccessfulFeeQuoteResponse),
		)

		httpmock.RegisterResponder("GET", policyUrl(urls[1]),
			httpmock.NewStringResponder(http.StatusOK, secondSuccessfulFeeQuoteResponse),
		)

		// when
		result, err := broadcaster.GetPolicyQuote(context.Background())

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
		assert.Equal(t, int64(100000000), result[0].Policy.MaxScriptSizePolicy)
	})

	t.Run("Should return error if all ArcClients return errors", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(2)

		httpmock.RegisterResponder("GET", policyUrl(urls[0]),
			httpmock.NewStringResponder(409, errorResponse409),
		)

		httpmock.RegisterResponder("GET", policyUrl(urls[1]),
			httpmock.NewStringResponder(409, errorResponse409),
		)

		// when
		result, err := broadcaster.GetPolicyQuote(context.Background())

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
		assert.EqualError(t, err, broadcast.ErrNoMinerResponse.Error())
	})

	t.Run("Should successfully query from single ArcClient", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(1)

		httpmock.RegisterResponder("GET", policyUrl(urls[0]),
			httpmock.NewStringResponder(http.StatusOK, firstSuccessfulFeeQuoteResponse),
		)

		// when
		result, err := broadcaster.GetPolicyQuote(context.Background())

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("Should return error if single ArcClient returns error", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(1)

		httpmock.RegisterResponder("GET", policyUrl(urls[0]),
			httpmock.NewStringResponder(409, errorResponse409),
		)

		// when
		result, err := broadcaster.GetPolicyQuote(context.Background())

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
