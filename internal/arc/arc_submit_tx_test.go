package arc

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/go-broadcast-client/common"
	"github.com/bitcoin-sv/go-broadcast-client/config"
	"github.com/bitcoin-sv/go-broadcast-client/internal/httpclient"
	"github.com/bitcoin-sv/go-broadcast-client/models"
	"github.com/bitcoin-sv/go-broadcast-client/shared"
	"github.com/stretchr/testify/assert"
)

func TestSubmitTransaction(t *testing.T) {
	testCases := []struct {
		name           string
		transaction    *common.Transaction
		httpResponse   *http.Response
		httpError      error
		expectedResult *models.SubmitTxResponse
		expectedError  error
	}{
		{
			name: "successful request",
			transaction: &common.Transaction{
				RawTx: "abc123",
			},
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
                    {
                        "txStatus": "CONFIRMED"
                    }
                    `)),
			},
			expectedResult: &models.SubmitTxResponse{
				TxStatus: common.Confirmed,
			},
		},
		{
			name: "error in HTTP request",
			transaction: &common.Transaction{
				RawTx: "abc123",
			},
			httpError:     errors.New("some error"),
			expectedError: errors.New("some error"),
		},
		{
			name: "missing txStatus in response",
			transaction: &common.Transaction{
				RawTx: "abc123",
			},
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
                    {
                        "dummyField": "dummyValue"
                    }
                    `)),
			},
			expectedError: shared.ErrMissingStatus,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			mockHttpClient := new(MockHttpClient)

			// create expected payload
			body, _ := createSubmitTxBody(tc.transaction)
			expectedPayload := httpclient.NewPayload(httpclient.POST, "http://example.com"+config.ArcSubmitTxRoute, "someToken", body)
			appendSubmitTxHeaders(&expectedPayload, tc.transaction)

			// define behavior of the mock
			mockHttpClient.On("DoRequest", context.Background(), expectedPayload).
				Return(tc.httpResponse, tc.httpError).Once()

			// use mock in the arc client
			client := &ArcClient{
				HTTPClient: mockHttpClient,
				apiURL:     "http://example.com",
				token:      "someToken",
			}

			// When
			result, err := client.SubmitTransaction(context.Background(), tc.transaction)

			// Then
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)

			// Assert Expectations on the mock
			mockHttpClient.AssertExpectations(t)
		})
	}
}
