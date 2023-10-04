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
				Hex: "abc123",
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
				Hex: "abc123",
			},
			httpError:     errors.New("some error"),
			expectedError: errors.New("some error"),
		},
		{
			name: "missing txStatus in response",
			transaction: &broadcast.Transaction{
				Hex: "abc123",
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

func TestConvertTransaction(t *testing.T) {
	testCases := []struct {
		name           string
		transaction    *broadcast.Transaction
		expectedResult *broadcast.Transaction
	}{
		{
			name: "successful conversion from RawTx to EF",
			transaction: &broadcast.Transaction{
				Hex: "0100000001a96fe5db0c2108e70abfb2b98ffbf4b7f66ca2a341c97484f1b1ebf967a2f51b000000006a47304402201846783d9e0e7abcaf3554b130f2e336865d67cbf18c5ad55580164a0b2a23590220614af1d8de08ffcbe3705de1fc48bd54031449cf8dba653da1af463922a6618d412102f9b7ecb5a0393e91aed5d27e35e723cf08c02979a8b0b1777c231a80b3d78d60ffffffff0232000000000000001976a914c63808bda42320a5a2425b3247e85cfc29f5e6f688ac31000000000000001976a9141d873677dfe9f3ae987c64fa3cb351194c68599988ac00000000",
			},
			expectedResult: &broadcast.Transaction{
				Hex: "010000000000000000ef01a96fe5db0c2108e70abfb2b98ffbf4b7f66ca2a341c97484f1b1ebf967a2f51b000000006a47304402201846783d9e0e7abcaf3554b130f2e336865d67cbf18c5ad55580164a0b2a23590220614af1d8de08ffcbe3705de1fc48bd54031449cf8dba653da1af463922a6618d412102f9b7ecb5a0393e91aed5d27e35e723cf08c02979a8b0b1777c231a80b3d78d60ffffffff64000000000000001976a91421c463658fc4457b937a7bb6aabd9c09fc70fcbb88ac0232000000000000001976a914c63808bda42320a5a2425b3247e85cfc29f5e6f688ac31000000000000001976a9141d873677dfe9f3ae987c64fa3cb351194c68599988ac00000000",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			test_tx := *tc.transaction

			// when
			err := convertToEfTransaction(&test_tx)

			// then
			assert.NoError(t, err)
			assert.Equal(t, *tc.expectedResult, test_tx)
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
				{Hex: "abc123"},
				{Hex: "cba321"},
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
				{Hex: "abc123"},
				{Hex: "cba321"},
			},
			httpError:     errors.New("some error"),
			expectedError: errors.New("some error"),
		},
		{
			name: "missing txStatus in response",
			transactions: []*broadcast.Transaction{
				{Hex: "abc123"},
				{Hex: "cba321"},
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
