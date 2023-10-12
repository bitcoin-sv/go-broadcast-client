package arc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/httpclient"
)

func TestSubmitTransaction(t *testing.T) {
	testCases := []struct {
		name           string
		transaction    *broadcast.Transaction
		httpResponse   *http.Response
		httpError      error
		expectedResult *broadcast.SubmitTxResponse
		expectedError  error
	}{
		{
			name: "successful request",
			transaction: &broadcast.Transaction{
				RawTx: "abc123",
			},
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`
                    {
                        "txStatus": "CONFIRMED"
                    }
                    `)),
			},
			expectedResult: &broadcast.SubmitTxResponse{
				BaseResponse: broadcast.BaseResponse{Miner: "http://example.com"},
				SubmittedTx:  &broadcast.SubmittedTx{TxStatus: broadcast.Confirmed},
			},
		},
		{
			name: "error in HTTP request",
			transaction: &broadcast.Transaction{
				RawTx: "abc123",
			},
			httpError:     errors.New("some error"),
			expectedError: errors.New("some error"),
		},
		{
			name: "missing txStatus in response",
			transaction: &broadcast.Transaction{
				RawTx: "abc123",
			},
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`
                    {
                        "dummyField": "dummyValue"
                    }
                    `)),
			},
			expectedError: broadcast.ErrMissingStatus,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			mockHttpClient := new(MockHttpClient)

			body, _ := createSubmitTxBody(tc.transaction)
			expectedPayload := httpclient.NewPayload(
				httpclient.POST,
				"http://example.com"+arcSubmitTxRoute,
				"someToken",
				body,
			)
			appendSubmitTxHeaders(&expectedPayload, nil)

			mockHttpClient.On("DoRequest", context.Background(), expectedPayload).
				Return(tc.httpResponse, tc.httpError).Once()

			client := &ArcClient{
				HTTPClient: mockHttpClient,
				apiURL:     "http://example.com",
				token:      "someToken",
			}

			// when
			result, err := client.SubmitTransaction(context.Background(), tc.transaction)

			// then
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)

			// assert Expectations on the mock
			mockHttpClient.AssertExpectations(t)
		})
	}
}

func TestSubmitBatchTransactions(t *testing.T) {
	testCases := []struct {
		name           string
		transactions   []*broadcast.Transaction
		httpResponse   *http.Response
		httpError      error
		expectedResult *broadcast.SubmitBatchTxResponse
		expectedError  error
	}{
		{
			name: "successful request",
			transactions: []*broadcast.Transaction{
				{RawTx: "abc123"},
				{RawTx: "cba321"},
			},
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`
					[
						{
							"txStatus": "CONFIRMED"
						},
						{
							"txStatus": "CONFIRMED"
						}
					]`)),
			},
			expectedResult: &broadcast.SubmitBatchTxResponse{
				BaseResponse: broadcast.BaseResponse{Miner: "http://example.com"},
				Transactions: []*broadcast.SubmittedTx{
					{TxStatus: broadcast.Confirmed},
					{TxStatus: broadcast.Confirmed},
				},
			},
		},
		{
			name: "error in HTTP request",
			transactions: []*broadcast.Transaction{
				{RawTx: "abc123"},
				{RawTx: "cba321"},
			},
			httpError:     errors.New("some error"),
			expectedError: errors.New("some error"),
		},
		{
			name: "missing txStatus in response",
			transactions: []*broadcast.Transaction{
				{RawTx: "abc123"},
				{RawTx: "cba321"},
			},
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`
					[
						{
							"dummyField": "dummyValue"
						},
						{
							"txStatus": "CONFIRMED"
						}
					]`)),
			},
			expectedError: broadcast.ErrMissingStatus,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			mockHttpClient := new(MockHttpClient)

			body, _ := createSubmitBatchTxsBody(tc.transactions)
			expectedPayload := httpclient.NewPayload(
				httpclient.POST,
				"http://example.com"+arcSubmitBatchTxsRoute,
				"someToken",
				body,
			)
			appendSubmitTxHeaders(&expectedPayload, nil)

			mockHttpClient.On("DoRequest", context.Background(), expectedPayload).
				Return(tc.httpResponse, tc.httpError).Once()

			client := &ArcClient{
				HTTPClient: mockHttpClient,
				apiURL:     "http://example.com",
				token:      "someToken",
			}

			// when
			result, err := client.SubmitBatchTransactions(context.Background(), tc.transactions)

			// then
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)

			// assert Expectations on the mock
			mockHttpClient.AssertExpectations(t)
		})
	}
}
