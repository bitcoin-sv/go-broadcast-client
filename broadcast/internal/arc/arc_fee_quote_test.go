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
	"github.com/stretchr/testify/mock"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

func TestFeeQuote(t *testing.T) {
	testCases := []struct {
		name           string
		httpResponse   *http.Response
		httpError      error
		expectedResult []*broadcast.FeeQuote
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
			expectedResult: []*broadcast.FeeQuote{
				{
					BaseResponse: broadcast.BaseResponse{Miner: "http://example.com"},
					MiningFee: broadcast.MiningFeeResponse{
						Bytes:    1000,
						Satoshis: 1,
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
			result, err := client.GetFeeQuote(context.Background())

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
