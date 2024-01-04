package arc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/rs/zerolog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock/fixtures"
)

func TestQueryTransaction(t *testing.T) {
	testCases := []struct {
		name           string
		httpResponse   *http.Response
		httpError      error
		expectedResult *broadcast.QueryTxResponse
		expectedError  error
	}{
		{
			name: "successful request",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`
					{
						"blockHash": "abc123",
						"txStatus": "CONFIRMED",
						"txid": "abc123"
					}
					`)),
			},
			expectedResult: &broadcast.QueryTxResponse{
				BaseResponse: broadcast.BaseResponse{Miner: "http://example.com"},
				BaseTxResponse: broadcast.BaseTxResponse{
					BlockHash: "abc123",
					TxStatus:  broadcast.Confirmed,
					TxID:      "abc123",
				},
			},
		},
		{
			name:          "error in HTTP request",
			httpError:     errors.New("some error"),
			expectedError: errors.New("some error"),
		},
		{
			name: "missing txID in response",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`
					{
						"blockHash": "abc123"
					}
					`)),
			},
			expectedError: ErrMissingTxID,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			mockHttpClient := new(MockHttpClient)
			testLogger := zerolog.Nop()

			mockHttpClient.On("DoRequest", context.Background(), mock.Anything).
				Return(tc.httpResponse, tc.httpError).Once()

			client := &ArcClient{
				HTTPClient: mockHttpClient,
				apiURL:     "http://example.com",
				token:      "someToken",
				Logger:     &testLogger,
			}

			// when
			result, err := client.QueryTransaction(context.Background(), "abc123")

			// then
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)

			// assert Expectations on the mock
			mockHttpClient.AssertExpectations(t)
		})
	}
}

func TestDecodeQueryResponseBody(t *testing.T) {
	testCases := []struct {
		name           string
		httpResponse   *http.Response
		expectedResult *broadcast.QueryTxResponse
	}{
		{
			name: "successful decode",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(fmt.Sprintf("{\"merklePath\":\"%s\"}", fixtures.TxMerklePath))),
			},
			expectedResult: &broadcast.QueryTxResponse{
				BaseResponse: broadcast.BaseResponse{Miner: "http://example.com"},
				BaseTxResponse: broadcast.BaseTxResponse{
					MerklePath: fixtures.TxMerklePath,
				},
			},
		},
		{
			name: "empty merkle path",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`
					{
						"merklePath": ""
					}
					`)),
			},
			expectedResult: &broadcast.QueryTxResponse{
				BaseResponse: broadcast.BaseResponse{
					Miner: "http://example.com",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			client := &ArcClient{
				apiURL: "http://example.com",
			}

			// when
			model, err := decodeQueryResponseBody(tc.httpResponse, client)

			// then
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, model)
		})
	}
}
