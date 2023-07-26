package arc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/go-broadcast-client/internal/httpclient"
	"github.com/bitcoin-sv/go-broadcast-client/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) DoRequest(method httpclient.HttpMethod, url string, token string, body io.Reader) (*http.Response, error) {
	args := m.Called(method, url, token, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestQueryTransaction(t *testing.T) {
	testCases := []struct {
		name           string
		httpResponse   *http.Response
		httpError      error
		expectedResult *models.QueryTxResponse
		expectedError  error
	}{
		{
			name: "successful request",
			httpResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
					{
						"blockHash": "abc123",
						"txStatus": "confirmed"
					}
					`)),
			},
			expectedResult: &models.QueryTxResponse{
				BlockHash: "abc123",
				TxStatus:  "confirmed",
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
						"txStatus": "confirmed"
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
			expectedError: ErrMissingStatus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockHttpClient := new(MockHttpClient)

			// define behavior of the mock
			mockHttpClient.On("DoRequest", httpclient.GET, mock.AnythingOfType("string"), mock.AnythingOfType("string"), nil).
				Return(tc.httpResponse, tc.httpError).Once()

			// use mock in the arc client
			client := &ArcClient{
				httpClient: mockHttpClient,
				apiURL:     "http://example.com",
				token:      "someToken",
			}

			result, err := client.QueryTransaction(context.Background(), "abc123")

			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)

			// Assert Expectations on the mock
			mockHttpClient.AssertExpectations(t)
		})
	}
}
