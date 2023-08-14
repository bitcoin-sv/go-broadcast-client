package arc

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	arcutils "github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc/utils"
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

	model := broadcast.QueryTxResponse{}
	err = arcutils.DecodeResponseBody(resp.Body, &model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func validateQueryTxResponse(model *broadcast.QueryTxResponse) error {
	if model.TxID == "" {
		return ErrMissingTxID
	}

	return nil
}
