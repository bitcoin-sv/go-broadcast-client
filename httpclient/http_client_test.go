package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

type mockModel struct {
	Str    string `json:"str"`
	Number int    `json:"number"`
}

func newMockModel(str string, number int) *mockModel {
	return &mockModel{
		Str:    str,
		Number: number,
	}
}

func modelJson(model *mockModel) string {
	json, _ := json.Marshal(model)
	return string(json)
}

func decodeModel(resp *http.Response) (*mockModel, error) {
	model := &mockModel{}
	err := json.NewDecoder(resp.Body).Decode(model)
	if err != nil {
		return nil, errors.New("Unable to decode response")
	}
	return model, err
}

const mockUrl = "http://test.com"

func nilErrorParser(statusCode int, body []byte) error {
	return nil
}

type mockError struct {
	ExtraInfo string `json:"extraInfo"`
}

func (e mockError) Error() string {
	return e.ExtraInfo
}

func mockErrorParser(statusCode int, body []byte) error {
	resultError := &mockError{}
	reader := bytes.NewReader(body)
	err := json.NewDecoder(reader).Decode(resultError)

	if err != nil {
		return errors.New("unable to decode an error")
	}

	return resultError
}

const mockErrorResponse = `
{
	"extraInfo": "mockErrorResponse"
}
`

func Test_HttpClient_RequestModel(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Should return model", func(t *testing.T) {
		httpmock.Reset()

		mockModel := newMockModel("test", 1)
		mockModelJson := modelJson(mockModel)

		httpmock.RegisterResponder("GET", mockUrl,
			httpmock.NewStringResponder(200, mockModelJson),
		)

		client := NewHttpClient()

		model, err := RequestModel(
			context.Background(),
			client.DoRequest,
			NewPayload(GET, mockUrl, "", nil),
			decodeModel,
			nilErrorParser,
		)

		assert.NoError(t, err)
		assert.Equal(t, mockModel, model)
	})

	t.Run("Should error if decode error", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder("GET", mockUrl,
			httpmock.NewStringResponder(200, "no-json"),
		)

		client := NewHttpClient()

		model, err := RequestModel(
			context.Background(),
			client.DoRequest,
			NewPayload(GET, mockUrl, "", nil),
			decodeModel,
			nilErrorParser,
		)

		assert.Error(t, err)
		assert.Nil(t, model)
	})

	t.Run("Should error if empty URL", func(t *testing.T) {
		httpmock.Reset()

		client := NewHttpClient()

		model, err := RequestModel(
			context.Background(),
			client.DoRequest,
			NewPayload(GET, "", "", nil),
			decodeModel,
			nilErrorParser,
		)

		assert.Error(t, err)
		assert.Nil(t, model)
	})

	t.Run("No success response", func(t *testing.T) {
		httpmock.Reset()

		client := NewHttpClient()

		httpmock.RegisterResponder("GET", mockUrl,
			httpmock.NewStringResponder(404, "ERROR"),
		)

		model, err := RequestModel(
			context.Background(),
			client.DoRequest,
			NewPayload(GET, mockUrl, "", nil),
			decodeModel,
			nilErrorParser,
		)

		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "404", "ERROR")
	})

	t.Run("No success response - special error", func(t *testing.T) {
		httpmock.Reset()

		client := NewHttpClient()

		httpmock.RegisterResponder("GET", mockUrl,
			httpmock.NewStringResponder(400, mockErrorResponse),
		)

		model, err := RequestModel(
			context.Background(),
			client.DoRequest,
			NewPayload(GET, mockUrl, "", nil),
			decodeModel,
			mockErrorParser,
		)

		assert.Error(t, err)
		assert.Nil(t, model)
		specialError, ok := err.(*mockError)
		assert.True(t, ok)
		assert.Equal(t, "mockErrorResponse", specialError.ExtraInfo)
	})
}
