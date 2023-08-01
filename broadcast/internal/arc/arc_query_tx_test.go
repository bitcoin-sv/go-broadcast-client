package arc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
						"txStatus": "CONFIRMED"
					}
					`)),
			},
			expectedResult: &broadcast.QueryTxResponse{
				BlockHash: "abc123",
				TxStatus:  broadcast.Confirmed,
			},
		},
		{
			name:          "error in HTTP request",
			httpError:     errors.New("some error"),
			expectedError: errors.New("some error"),
		},
		{
			name: "missing blockHash in response",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`
					{
						"txStatus": "CONFIRMED"
					}
					`)),
			},
			expectedError: ErrMissingHash,
		},
		{
			name: "missing txStatus in response",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`
					{
						"blockHash": "abc123"
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

			mockHttpClient.On("DoRequest", context.Background(), mock.Anything).
				Return(tc.httpResponse, tc.httpError).Once()

			client := &ArcClient{
				HTTPClient: mockHttpClient,
				apiURL:     "http://example.com",
				token:      "someToken",
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
