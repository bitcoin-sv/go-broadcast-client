package arc

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-api"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/httpclient"
)

var ErrMissingHash = errors.New("missing tx hash")

func (a *ArcClient) QueryTransaction(ctx context.Context, txID string) (*broadcast_api.QueryTxResponse, error) {
	if a == nil {
		return nil, broadcast_api.ErrClientUndefined
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
func queryTransaction(ctx context.Context, arc *ArcClient, txHash string) (*broadcast_api.QueryTxResponse, error) {
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

func decodeQueryTxBody(body io.ReadCloser) (*broadcast_api.QueryTxResponse, error) {
	model := broadcast_api.QueryTxResponse{}
	err := json.NewDecoder(body).Decode(&model)

	if err != nil {
		return nil, broadcast_api.ErrUnableToDecodeResponse
	}

	return &model, nil
}

func validateQueryTxResponse(model *broadcast_api.QueryTxResponse) error {

	if model.BlockHash == "" {
		return ErrMissingHash
	}

	if model.TxStatus == "" {
		return broadcast_api.ErrMissingStatus
	}

	return nil
}
