package acceptancetests

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/stretchr/testify/assert"
)

const successfulSubmitResponse = `
{
    "txStatus": "MINED"
}
`

const successfulSubmitBatchResponse = `
[
	{
    	"txStatus": "MINED"
	},
	{
    	"txStatus": "MINED"
	}
]
`

const missingStatusSubmitResponse = `
{
	"blockHash": "hash"
}
`

const mockHex = "transaction-data"

func mockBatch() []*broadcast.Transaction {
	return []*broadcast.Transaction{
		{Hex: "transaction-0-data"},
		{Hex: "transaction-1-data"},
	}
}

func submitUrl(base string) string {
	return requestUrl(base, "/v1/tx")
}

func submitBatchUrl(base string) string {
	return requestUrl(base, "/v1/txs")
}

func TestSubmitTransaction(t *testing.T) {
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
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())

		assert.Equal(t, 1, callCount(submitUrl(urls[0])))
		assert.Equal(t, 0, callCount(submitUrl(urls[1])))
	})

	t.Run("Should return error if both ArcClients return errors", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(2)

		httpmock.RegisterResponder("POST", submitUrl(urls[0]),
			httpmock.NewStringResponder(409, errorResponse409),
		)

		httpmock.RegisterResponder("POST", submitUrl(urls[1]),
			httpmock.NewStringResponder(409, errorResponse409),
		)

		// when
		result, err := broadcaster.SubmitTransaction(context.Background(), &broadcast.Transaction{Hex: "transaction-data"})

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, broadcast.ErrAllBroadcastersFailed.Error())
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
		assert.Equal(t, 1, callCount(submitUrl(urls[0])))
		assert.Equal(t, 1, callCount(submitUrl(urls[1])))
	})

	t.Run("Should return error if one ArcClient returns error and the other returns invalid response", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(2)

		httpmock.RegisterResponder("POST", submitUrl(urls[0]),
			httpmock.NewStringResponder(409, errorResponse409),
		)

		httpmock.RegisterResponder("POST", submitUrl(urls[1]),
			httpmock.NewErrorResponder(errors.New("http error")),
		)

		// when
		result, err := broadcaster.SubmitTransaction(context.Background(), &broadcast.Transaction{Hex: "transaction-data"})

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
	})

	t.Run("Should successfully submit transaction if one ArcClient returns missing status and the other is successful", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(2)

		httpmock.RegisterResponder("POST", submitUrl(urls[0]),
			httpmock.NewStringResponder(http.StatusOK, missingStatusSubmitResponse),
		)

		httpmock.RegisterResponder("POST", submitUrl(urls[1]),
			httpmock.NewStringResponder(http.StatusOK, successfulSubmitResponse),
		)

		// when
		result, err := broadcaster.SubmitTransaction(context.Background(), &broadcast.Transaction{Hex: mockHex})

		// then
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
		assert.Equal(t, urls[1], result.Miner)
	})
}

func TestSubmitBatchTransactions(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Should successfully submit batch of transactions using first of few ArcClients", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(2)

		httpmock.RegisterResponder("POST", submitBatchUrl(urls[0]),
			httpmock.NewStringResponder(http.StatusOK, successfulSubmitBatchResponse),
		)

		httpmock.RegisterResponder("POST", submitBatchUrl(urls[1]),
			httpmock.NewStringResponder(http.StatusOK, successfulSubmitBatchResponse),
		)

		batch := mockBatch()

		// when
		result, err := broadcaster.SubmitBatchTransactions(context.Background(), batch)

		// then
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
		assert.Equal(t, 1, callCount(submitBatchUrl(urls[0])))
	})

	t.Run("Should successfully submit batch of transactions if one ArcClient returns missing status and the other is successful", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(2)

		httpmock.RegisterResponder("POST", submitBatchUrl(urls[0]),
			httpmock.NewStringResponder(http.StatusOK, `[{"blockHash": "hash"}, {"txStatus": "MINED"}`),
		)

		httpmock.RegisterResponder("POST", submitBatchUrl(urls[1]),
			httpmock.NewStringResponder(http.StatusOK, successfulSubmitBatchResponse),
		)

		batch := mockBatch()

		// when
		result, err := broadcaster.SubmitBatchTransactions(context.Background(), batch)

		// then
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
		assert.Equal(t, 1, callCount(submitBatchUrl(urls[1])))
		assert.Equal(t, urls[1], result.Miner)
	})

	t.Run("Should return error if every ArcClients return errors", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(2)

		httpmock.RegisterResponder("POST", submitBatchUrl(urls[0]),
			httpmock.NewErrorResponder(errors.New("http error")),
		)

		httpmock.RegisterResponder("POST", submitBatchUrl(urls[1]),
			httpmock.NewErrorResponder(errors.New("http error")),
		)

		batch := mockBatch()

		// when
		result, err := broadcaster.SubmitBatchTransactions(context.Background(), batch)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, broadcast.ErrAllBroadcastersFailed.Error())
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
	})

	t.Run("Should return error if one ArcClient returns error and the other returns invalid response", func(t *testing.T) {
		httpmock.Reset()
		// given
		broadcaster, urls := getBroadcaster(2)

		httpmock.RegisterResponder("POST", submitBatchUrl(urls[0]),
			httpmock.NewStringResponder(409, errorResponse409),
		)

		httpmock.RegisterResponder("POST", submitBatchUrl(urls[1]),
			httpmock.NewErrorResponder(errors.New("http error")),
		)

		batch := mockBatch()

		// when
		result, err := broadcaster.SubmitBatchTransactions(context.Background(), batch)

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
	})
}
