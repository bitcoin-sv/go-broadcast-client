package arc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
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
					BaseTxResponse: broadcast.BaseTxResponse{
						TxStatus: broadcast.Confirmed,
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
			appendSubmitTxHeaders(&expectedPayload, nil, client.deploymentID)

			mockHttpClient.On("DoRequest", context.Background(), expectedPayload).
				Return(tc.httpResponse, tc.httpError).Once()

			// when
			result, err := client.SubmitTransaction(context.Background(), tc.transaction)

			// then
			assert.Equal(t, tc.expectedResult, result)

			if err != nil {
				assert.True(t, strings.Contains(err.Error(), tc.expectedError.Error()))
			}
			// assert Expectations on the mock
			mockHttpClient.AssertExpectations(t)
		})
	}
}

func TestFormatTxRequest(t *testing.T) {
	testCases := []struct {
		name           string
		httpResponse   *http.Response
		httpError      error
		expectedResult *SubmitTxRequest
		expectedError  error
	}{
		{
			name: "successful request",
		},
		{
			name: "error in HTTP request",
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

			rawTx := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff1703db0a0d2f43555656452f8313c10017f13d0fe1f00800ffffffff01d140a212000000001976a914d648686cf603c11850f39600e37312738accca8f88ac00000000"

			// when
			result, err := rawTxRequest(client, rawTx)

			// t.Logf(result.RawTx)

			// then
			assert.Equal(t, tc.expectedResult, result)

			if err != nil {
				assert.True(t, strings.Contains(err.Error(), tc.expectedError.Error()))
			}

			// assert Expectations on the mock
			mockHttpClient.AssertExpectations(t)
		})
	}
}
