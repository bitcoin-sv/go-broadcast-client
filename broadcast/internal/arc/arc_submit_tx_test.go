package arc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-api"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/httpclient"
	"github.com/stretchr/testify/assert"
)

func TestSubmitTransaction(t *testing.T) {
	testCases := []struct {
		name           string
		transaction    *broadcast_api.Transaction
		httpResponse   *http.Response
		httpError      error
		expectedResult *broadcast_api.SubmitTxResponse
		expectedError  error
	}{
		{
			name: "successful request",
			transaction: &broadcast_api.Transaction{
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
			expectedResult: &broadcast_api.SubmitTxResponse{
				TxStatus: broadcast_api.Confirmed,
			},
		},
		{
			name: "error in HTTP request",
			transaction: &broadcast_api.Transaction{
				RawTx: "abc123",
			},
			httpError:     errors.New("some error"),
			expectedError: errors.New("some error"),
		},
		{
			name: "missing txStatus in response",
			transaction: &broadcast_api.Transaction{
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
			expectedError: broadcast_api.ErrMissingStatus,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			mockHttpClient := new(MockHttpClient)

			body, _ := createSubmitTxBody(tc.transaction)
			expectedPayload := httpclient.NewPayload(httpclient.POST, "http://example.com"+arcSubmitTxRoute, "someToken", body)
			appendSubmitTxHeaders(&expectedPayload, tc.transaction)

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
