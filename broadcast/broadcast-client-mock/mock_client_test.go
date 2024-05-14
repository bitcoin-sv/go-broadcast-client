package broadcast_client_mock

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock/fixtures"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc/mocks"
)

func TestMockClientSuccess(t *testing.T) {
	t.Run("Should successfully query for Policy Quote from mock Arc Client with Success Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockSuccess).
			Build()
		expectedResult := []*broadcast.PolicyQuoteResponse{mocks.Policy1, mocks.Policy2}

		// when
		result, fail := broadcaster.GetPolicyQuote(context.Background())

		// then
		assert.NoError(t, fail)
		assert.NotNil(t, result)
		assert.Equal(t, result, expectedResult)
	})

	t.Run("Should successfully query for Fee Quote from mock Arc Client with Success Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockSuccess).
			Build()
		expectedResult := []*broadcast.FeeQuote{mocks.Fee1, mocks.Fee2}

		// when
		result, fail := broadcaster.GetFeeQuote(context.Background())

		// then
		assert.Nil(t, fail)
		assert.NotNil(t, result)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Should successfully query for transaction from mock Arc Client with Success Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockSuccess).
			Build()
		testTxId := "test-txid"

		// when
		result, err := broadcaster.QueryTransaction(context.Background(), testTxId)

		// then
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, result.Miner, fixtures.ProviderMain)
		assert.Equal(t, result.TxID, testTxId)
		assert.Equal(t, result.BlockHash, fixtures.TxBlockHash)
		assert.Equal(t, result.BlockHeight, fixtures.TxBlockHeight)
	})

	t.Run("Should return successful submit transaction response from mock Arc Client with Success Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockSuccess).
			Build()

		// when
		result, fail := broadcaster.SubmitTransaction(context.Background(), &broadcast.Transaction{Hex: "test-rawtx"})

		// then
		assert.Nil(t, fail)
		assert.NotNil(t, result)
		assert.Equal(t, result.Miner, fixtures.ProviderMain)
		assert.Equal(t, result.BlockHash, fixtures.TxBlockHash)
		assert.Equal(t, result.BlockHeight, fixtures.TxBlockHeight)
		assert.Equal(t, result.TxStatus, fixtures.TxStatus)
		assert.Equal(t, result.Status, fixtures.TxResponseStatus)
		assert.Equal(t, result.Title, fixtures.TxResponseTitle)
	})

	t.Run("Should return successful submit batch transactions response from mock Arc Client with Success Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockSuccess).
			Build()
		expectedResult := &broadcast.SubmitBatchTxResponse{
			BaseResponse: broadcast.BaseResponse{Miner: fixtures.ProviderMain},
			Transactions: []*broadcast.SubmittedTx{
				mocks.SubmittedTx,
				mocks.SubmittedTxSecondary,
			},
		}

		// when
		result, fail := broadcaster.SubmitBatchTransactions(context.Background(), []*broadcast.Transaction{{Hex: "test-rawtx"}, {Hex: "test2-rawtx"}})

		// then
		assert.Nil(t, fail)
		assert.NotNil(t, result)
		assert.Equal(t, expectedResult, result)
	})
}

func TestMockClientFailure(t *testing.T) {
	t.Run("Should return error from GetPolicyQuote method of mock Arc Client with Failure Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockFailure).
			Build()

		// when
		result, fail := broadcaster.GetPolicyQuote(context.Background())

		// then
		assert.Error(t, fail)
		assert.Nil(t, result)
		assert.ErrorContains(t, fail, broadcast.ErrNoMinerResponse.Error())
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
		assert.ErrorContains(t, err, broadcast.ErrNoMinerResponse.Error())
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
		assert.ErrorContains(t, err, broadcast.ErrNoMinerResponse.Error())
	})

	t.Run("Should return error from SubmitTransaction method of mock Arc Client with Failure Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockFailure).
			Build()

		// when
		result, err := broadcaster.SubmitTransaction(context.Background(), &broadcast.Transaction{Hex: "test-rawtx"})

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, broadcast.ErrAllBroadcastersFailed.Error())
	})

	t.Run("Should return error from SubmitBatchTransaction method of mock Arc Client with Failure Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockFailure).
			Build()

		// when
		result, err := broadcaster.SubmitBatchTransactions(context.Background(), []*broadcast.Transaction{{Hex: "test-rawtx"}, {Hex: "test2-rawtx"}})

		// then
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorContains(t, err, broadcast.ErrAllBroadcastersFailed.Error())
	})
}

func TestMockClientTimeout(t *testing.T) {
	const defaultTestTime = 200 * time.Millisecond

	t.Run("Should successfully query for Policy Quote after a timeout period from mock Arc Client with Timeout Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockTimeout).
			Build()
		ctx, cancel := context.WithTimeout(context.Background(), defaultTestTime)
		defer cancel()
		startTime := time.Now()
		expectedResult := []*broadcast.PolicyQuoteResponse{mocks.Policy1, mocks.Policy2}

		// when
		result, err := broadcaster.GetPolicyQuote(ctx)

		// then
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Greater(t, time.Since(startTime), defaultTestTime)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Should successfully query for Fee Quote after a timeout period from mock Arc Client with Timeout Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockTimeout).
			Build()
		ctx, cancel := context.WithTimeout(context.Background(), defaultTestTime)
		defer cancel()
		startTime := time.Now()
		expectedResult := []*broadcast.FeeQuote{mocks.Fee1, mocks.Fee2}

		// when
		result, err := broadcaster.GetFeeQuote(ctx)

		// then
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Greater(t, time.Since(startTime), defaultTestTime)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Should successfully query for transaction after a timeout period from mock Arc Client with Timeout Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockTimeout).
			Build()
		ctx, cancel := context.WithTimeout(context.Background(), defaultTestTime)
		defer cancel()
		startTime := time.Now()
		testTxId := "test-txid"

		// when
		result, err := broadcaster.QueryTransaction(ctx, testTxId)

		// then
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Greater(t, time.Since(startTime), defaultTestTime)
		assert.Equal(t, result.Miner, fixtures.ProviderMain)
		assert.Equal(t, result.TxID, testTxId)
		assert.Equal(t, result.BlockHash, fixtures.TxBlockHash)
		assert.Equal(t, result.BlockHeight, fixtures.TxBlockHeight)
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
		result, err := broadcaster.SubmitTransaction(ctx, &broadcast.Transaction{Hex: "test-rawtx"})

		// then
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Greater(t, time.Since(startTime), defaultTestTime)
		assert.Equal(t, result.Miner, fixtures.ProviderMain)
		assert.Equal(t, result.BlockHash, fixtures.TxBlockHash)
		assert.Equal(t, result.BlockHeight, fixtures.TxBlockHeight)
		assert.Equal(t, result.TxStatus, fixtures.TxStatus)
		assert.Equal(t, result.Status, fixtures.TxResponseStatus)
		assert.Equal(t, result.Title, fixtures.TxResponseTitle)
	})

	t.Run("Should return successful submit batch transactions response after a timeout period from mock Arc Client with Timeout Mock Type", func(t *testing.T) {
		// given
		broadcaster := Builder().
			WithMockArc(MockTimeout).
			Build()
		ctx, cancel := context.WithTimeout(context.Background(), defaultTestTime)
		defer cancel()
		startTime := time.Now()
		expectedResult := &broadcast.SubmitBatchTxResponse{
			BaseResponse: broadcast.BaseResponse{Miner: fixtures.ProviderMain},
			Transactions: []*broadcast.SubmittedTx{
				mocks.SubmittedTx,
				mocks.SubmittedTxSecondary,
			},
		}

		// when
		result, err := broadcaster.SubmitBatchTransactions(ctx, []*broadcast.Transaction{{Hex: "test-rawtx"}, {Hex: "test2-rawtx"}})

		// then
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Greater(t, time.Since(startTime), defaultTestTime)
		assert.Equal(t, expectedResult, result)
	})
}
