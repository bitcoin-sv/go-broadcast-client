package arc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/rs/zerolog"

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
				SubmittedTx: &broadcast.SubmittedTx{
					BaseSubmitTxResponse: broadcast.BaseSubmitTxResponse{
						BaseTxResponse: broadcast.BaseTxResponse{
							TxStatus: broadcast.Confirmed,
						},
					},
				},
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
			testLogger := zerolog.Nop()

			client := &ArcClient{
				HTTPClient: mockHttpClient,
				apiURL:     "http://example.com",
				token:      "someToken",
				Logger:     &testLogger,
			}

			body, _ := createSubmitTxBody(client, tc.transaction, broadcast.EfFormat)
			expectedPayload := httpclient.NewPayload(
				httpclient.POST,
				"http://example.com"+arcSubmitTxRoute,
				"someToken",
				body,
			)
			appendSubmitTxHeaders(&expectedPayload, nil, client.headers)

			mockHttpClient.On("DoRequest", context.Background(), expectedPayload).
				Return(tc.httpResponse, tc.httpError).Once()

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
