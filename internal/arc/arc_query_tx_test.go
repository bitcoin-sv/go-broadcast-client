package arc

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/shared"
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
				Body: ioutil.NopCloser(bytes.NewBufferString(`
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
				Body: ioutil.NopCloser(bytes.NewBufferString(`
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
				Body: ioutil.NopCloser(bytes.NewBufferString(`
					{
						"blockHash": "abc123"
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

			// define behavior of the mock
			mockHttpClient.On("DoRequest", context.Background(), mock.Anything).
				Return(tc.httpResponse, tc.httpError).Once()

			// use mock in the arc client
			client := &ArcClient{
				HTTPClient: mockHttpClient,
				apiURL:     "http://example.com",
				token:      "someToken",
			}

			// When
			result, err := client.QueryTransaction(context.Background(), "abc123")

			// Then
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)

			// Assert Expectations on the mock
			mockHttpClient.AssertExpectations(t)
		})
	}
}
