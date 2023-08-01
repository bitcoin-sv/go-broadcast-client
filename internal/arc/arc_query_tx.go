package arc

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/go-broadcast-client/internal/httpclient"
	"github.com/bitcoin-sv/go-broadcast-client/shared"
)

var ErrMissingHash = errors.New("missing tx hash")

func (a *ArcClient) QueryTransaction(ctx context.Context, txID string) (*broadcast.QueryTxResponse, error) {
	if a == nil {
		return nil, shared.ErrClientUndefined
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
	url := arc.apiURL + broadcast.ArcQueryTxRoute + txHash
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
		return nil, shared.ErrUnableToDecodeResponse
	}

	return &model, nil
}

func validateQueryTxResponse(model *broadcast.QueryTxResponse) error {

	if model.BlockHash == "" {
		return ErrMissingHash
	}

	if model.TxStatus == "" {
		return shared.ErrMissingStatus
	}

	return nil
}
