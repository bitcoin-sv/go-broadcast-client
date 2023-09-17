package broadcast_client_mock

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

func TestMockClientSuccess(t *testing.T) {
	t.Run("Should successfully query for Policy Quote from mock Arc Client with Success Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockSuccess).
			Build()

		// when
		result, err := broadcaster.GetPolicyQuote(context.Background())

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Should successfully query for Fee Quote from mock Arc Client with Success Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockSuccess).
			Build()

		// when
		result, err := broadcaster.GetFeeQuote(context.Background())

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Should successfully query for transaction from mock Arc Client with Success Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockSuccess).
			Build()

		// when
		result, err := broadcaster.QueryTransaction(context.Background(), "test-txid")

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Should return successful submit transaction response from mock Arc Client with Success Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockSuccess).
			Build()

		// when
		result, err := broadcaster.SubmitTransaction(context.Background(), &broadcast.Transaction{RawTx: "test-rawtx"})

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Should return successful submit batch transactions response from mock Arc Client with Success Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockSuccess).
			Build()

		// when
		result, err := broadcaster.SubmitBatchTransactions(context.Background(), []*broadcast.Transaction{{RawTx: "test-rawtx"}, {RawTx: "test2-rawtx"}})

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestMockClientFailure(t *testing.T) {
	t.Run("Should return error from GetPolicyQuote method of mock Arc Client with Failure Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockFailure).
			Build()

		// when
		result, err := broadcaster.GetPolicyQuote(context.Background())

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, broadcast.ErrNoMinerResponse.Error())
	})

	t.Run("Should return error from GetFeeQuote method of mock Arc Client with Failure Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockFailure).
			Build()

		// when
		result, err := broadcaster.GetFeeQuote(context.Background())

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, broadcast.ErrNoMinerResponse.Error())
	})

	t.Run("Should return error from QueryTransaction method of mock Arc Client with Failure Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockFailure).
			Build()

		// when
		result, err := broadcaster.QueryTransaction(context.Background(), "test-txid")

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, broadcast.ErrAllBroadcastersFailed.Error())
	})

	t.Run("Should return error from SubmitTransaction method of mock Arc Client with Failure Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockFailure).
			Build()

		// when
		result, err := broadcaster.SubmitTransaction(context.Background(), &broadcast.Transaction{RawTx: "test-rawtx"})

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, broadcast.ErrAllBroadcastersFailed.Error())
	})

	t.Run("Should return error from SubmitBatchTransaction method of mock Arc Client with Failure Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockFailure).
			Build()

		// when
		result, err := broadcaster.SubmitBatchTransactions(context.Background(), []*broadcast.Transaction{{RawTx: "test-rawtx"}, {RawTx: "test2-rawtx"}})

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, broadcast.ErrAllBroadcastersFailed.Error())
	})
}

func TestMockClientTimeout(t *testing.T) {
	const defaultTestTime = 200*time.Millisecond

	t.Run("Should successfully query for Policy Quote after a timeout period from mock Arc Client with Timeout Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockTimeout).
			Build()
		ctx, cancel := context.WithTimeout(context.Background(), defaultTestTime)
		defer cancel()
		startTime := time.Now()

		// when
		result, err := broadcaster.GetPolicyQuote(ctx)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Greater(t, time.Now().Sub(startTime), defaultTestTime)
	})

	t.Run("Should successfully query for Fee Quote after a timeout period from mock Arc Client with Timeout Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockTimeout).
			Build()
		ctx, cancel := context.WithTimeout(context.Background(), defaultTestTime)
		defer cancel()
		startTime := time.Now()

		// when
		result, err := broadcaster.GetFeeQuote(ctx)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Greater(t, time.Now().Sub(startTime), defaultTestTime)
	})

	t.Run("Should successfully query for transaction after a timeout period from mock Arc Client with Timeout Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockTimeout).
			Build()
		ctx, cancel := context.WithTimeout(context.Background(), defaultTestTime)
		defer cancel()
		startTime := time.Now()

		// when
		result, err := broadcaster.QueryTransaction(ctx, "test-txid")

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Greater(t, time.Now().Sub(startTime), defaultTestTime)
	})

	t.Run("Should return successful submit transaction response after a timeout period from mock Arc Client with Timeout Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockTimeout).
			Build()
		ctx, cancel := context.WithTimeout(context.Background(), defaultTestTime)
		defer cancel()
		startTime := time.Now()

		// when
		result, err := broadcaster.SubmitTransaction(ctx, &broadcast.Transaction{RawTx: "test-rawtx"})

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Greater(t, time.Now().Sub(startTime), defaultTestTime)
	})

	t.Run("Should return successful submit batch transactions response after a timeout period from mock Arc Client with Timeout Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockTimeout).
			Build()
		ctx, cancel := context.WithTimeout(context.Background(), defaultTestTime)
		defer cancel()
		startTime := time.Now()

		// when
		result, err := broadcaster.SubmitBatchTransactions(ctx, []*broadcast.Transaction{{RawTx: "test-rawtx"}, {RawTx: "test2-rawtx"}})

		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Greater(t, time.Now().Sub(startTime), defaultTestTime)
	})
}
