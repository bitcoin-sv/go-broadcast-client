package acceptancetests

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/stretchr/testify/assert"
)

const firstSuccessfulTxResponse = `
{
    "blockHash": "hash",
    "txStatus": "MINED",
    "txid": "abc123"
}
`

const secondSuccessfulTxResponse = `
{
    "blockHash": "hash",
    "txStatus": "CONFIRMED",
    "txid": "abc123"
}
`

const mockTxID = "txID"

func txUrl(base string) string {
	return requestUrl(base, "/v1/tx/"+mockTxID)
}

func TestQueryTransaction(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Should successfully query when both can return success", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(2)

		httpmock.RegisterResponder("GET", txUrl(urls[0]),
			httpmock.NewStringResponder(http.StatusOK, firstSuccessfulTxResponse),
		)
		// first miner responded successfully, next one should be skipped
		httpmock.RegisterResponder("GET", txUrl(urls[1]),
			httpmock.NewStringResponder(http.StatusOK, secondSuccessfulTxResponse),
		)

		// when
		result, err := broadcaster.QueryTransaction(context.Background(), mockTxID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
		assert.Equal(t, "MINED", string(result.TxStatus))
	})

	t.Run("Should successfully query when first one can return success", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(2)

		httpmock.RegisterResponder("GET", txUrl(urls[0]),
			httpmock.NewStringResponder(http.StatusOK, firstSuccessfulTxResponse),
		)
		// first miner responded successfully, next one should be skipped
		httpmock.RegisterResponder("GET", txUrl(urls[1]),
			httpmock.NewStringResponder(409, errorResponse409),
		)

		// when
		result, err := broadcaster.QueryTransaction(context.Background(), mockTxID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
		assert.Equal(t, "MINED", string(result.TxStatus))
	})

	t.Run("Should successfully query when second one can return success", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(2)

		httpmock.RegisterResponder("GET", txUrl(urls[0]),
			httpmock.NewStringResponder(409, errorResponse409),
		)
		// this time, first miner responded with error, second one should be queried
		httpmock.RegisterResponder("GET", txUrl(urls[1]),
			httpmock.NewStringResponder(http.StatusOK, secondSuccessfulTxResponse),
		)

		// when
		result, err := broadcaster.QueryTransaction(context.Background(), mockTxID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
		assert.Equal(t, "CONFIRMED", string(result.TxStatus))
	})

	t.Run("Should return error if all ArcClients return errors", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(2)

		httpmock.RegisterResponder("GET", txUrl(urls[0]),
			httpmock.NewStringResponder(409, errorResponse409),
		)

		httpmock.RegisterResponder("GET", txUrl(urls[1]),
			httpmock.NewStringResponder(409, errorResponse409),
		)

		// when
		result, err := broadcaster.QueryTransaction(context.Background(), mockTxID)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, broadcast.ErrAllBroadcastersFailed.Error())
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
	})

	t.Run("Should successfully query from single ArcClient", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(2)

		httpmock.RegisterResponder("GET", txUrl(urls[0]),
			httpmock.NewStringResponder(http.StatusOK, firstSuccessfulTxResponse),
		)

		// when
		result, err := broadcaster.QueryTransaction(context.Background(), mockTxID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Should return error if single ArcClient returns error", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(1)

		httpmock.RegisterResponder("GET", txUrl(urls[0]),
			httpmock.NewStringResponder(409, errorResponse409),
		)

		// when
		result, err := broadcaster.QueryTransaction(context.Background(), mockTxID)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
