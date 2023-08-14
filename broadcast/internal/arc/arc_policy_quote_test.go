package arc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

func TestPolicyQuote(t *testing.T) {
	testCases := []struct {
		name           string
		httpResponse   *http.Response
		httpError      error
		expectedResult []*broadcast.PolicyQuoteResponse
		expectedError  error
	}{
		{
			name: "successful request",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBufferString(`
                    {
                      "policy": {
                        "maxscriptsizepolicy": 100000000,
                        "maxtxsigopscountspolicy": 4294967295,
                        "maxtxsizepolicy": 100000000,
                        "miningFee": {
                          "bytes": 1000,
                          "satoshis": 1
                        }
                      },
                      "timestamp": "2023-08-10T13:49:07.308687569Z"
                    }
					`)),
			},
			expectedResult: []*broadcast.PolicyQuoteResponse{
				{
					Miner: "http://example.com",
					Policy: broadcast.PolicyResponse{
						MaxScriptSizePolicy:    100000000,
						MaxTxSigOpsCountPolicy: 4294967295,
						MaxTxSizePolicy:        100000000,
						MiningFee: broadcast.MiningFeeResponse{
							Bytes:    1000,
							Satoshis: 1,
						},
					},
					Timestamp: "2023-08-10T13:49:07.308687569Z",
				},
			},
		},
		{
			name:          "error in HTTP request",
			httpError:     errors.New("some error"),
			expectedError: errors.New("some error"),
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
			result, err := client.GetPolicyQuote(context.Background())

			// then
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)

			// assert Expectations on the mock
			mockHttpClient.AssertExpectations(t)
		})
	}
}
