package arc

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/httpclient"
)

var ErrMissingTxID = errors.New("missing tx id")

func (a *ArcClient) QueryTransaction(ctx context.Context, txID string) (*broadcast.QueryTxResponse, error) {
	if a == nil {
		return nil, broadcast.ErrClientUndefined
	}

	result, err := queryTransaction(ctx, a, txID)
	if err != nil {
		return nil, err
	}

	err = validateQueryTxResponse(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// queryTransaction will fire the HTTP request to retrieve the tx status and details
func queryTransaction(ctx context.Context, arc *ArcClient, txHash string) (*broadcast.QueryTxResponse, error) {
	url := arc.apiURL + arcQueryTxRoute + txHash
	pld := httpclient.NewPayload(
		httpclient.GET,
		url,
		arc.token,
		nil,
	)

	resp, err := arc.HTTPClient.DoRequest(
		ctx,
		pld,
	)
	if err != nil {
		return nil, err
	}

	model, err := decodeQueryTxBody(resp.Body)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func decodeQueryTxBody(body io.ReadCloser) (*broadcast.QueryTxResponse, error) {
	model := broadcast.QueryTxResponse{}
	err := json.NewDecoder(body).Decode(&model)

	if err != nil {
		return nil, broadcast.ErrUnableToDecodeResponse
	}

	return &model, nil
}

func validateQueryTxResponse(model *broadcast.QueryTxResponse) error {
	if model.TxID == "" {
		return ErrMissingTxID
	}

	return nil
}
