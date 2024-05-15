package acceptancetests

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jarcoal/httpmock"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

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
			httpmock.NewStringResponder(http.StatusOK, firstSuccessfulPolicyResponse),
		)

		httpmock.RegisterResponder("GET", policyUrl(urls[1]),
			httpmock.NewStringResponder(http.StatusOK, secondSuccessfulPolicyResponse),
		)

		// when
		result, err := broadcaster.GetFeeQuote(context.Background())

		// then
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
		assert.Equal(t, int64(1), result[0].MiningFee.Satoshis)
		assert.Equal(t, urls[0], result[0].Miner)
		assert.Equal(t, int64(2), result[1].MiningFee.Satoshis)
		assert.Equal(t, urls[1], result[1].Miner)
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
		result, err := broadcaster.GetFeeQuote(context.Background())

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
		assert.Contains(t, err.Error(), broadcast.ErrNoMinerResponse.Error())
	})

	t.Run("Should successfully query from single ArcClient", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(1)

		httpmock.RegisterResponder("GET", policyUrl(urls[0]),
			httpmock.NewStringResponder(http.StatusOK, firstSuccessfulPolicyResponse),
		)

		// when
		result, err := broadcaster.GetFeeQuote(context.Background())

		// then
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
		assert.Equal(t, int64(1), result[0].MiningFee.Satoshis)
	})

	t.Run("Should return error if single ArcClient returns error", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(1)

		httpmock.RegisterResponder("GET", policyUrl(urls[0]),
			httpmock.NewStringResponder(409, errorResponse409),
		)

		// when
		result, err := broadcaster.GetFeeQuote(context.Background())

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
