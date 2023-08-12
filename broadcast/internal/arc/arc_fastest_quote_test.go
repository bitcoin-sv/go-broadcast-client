package arc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

func TestFastestQuote(t *testing.T) {
	testCases := []struct {
		name           string
		httpResponse   *http.Response
		httpError      error
		expectedResult *broadcast.FeeQuote
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
			expectedResult: &broadcast.FeeQuote{
				Miner: "http://example.com",
				MiningFee: broadcast.MiningFeeResponse{
					Bytes:    1000,
					Satoshis: 1,
				},
				Timestamp: "2023-08-10T13:49:07.308687569Z",
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

			mockHttpClient.On("DoRequest", mock.Anything, mock.Anything).
				Return(tc.httpResponse, tc.httpError).Once()

			client := &ArcClient{
				HTTPClient: mockHttpClient,
				apiURL:     "http://example.com",
				token:      "someToken",
			}

			// when
            result, err := client.GetFastestQuote(context.Background(), time.Second)

			// then
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)

			// assert Expectations on the mock
			mockHttpClient.AssertExpectations(t)
		})
	}
}