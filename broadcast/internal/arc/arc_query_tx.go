package arc

import (
	"context"
	"errors"
	"net/http"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	arcutils "github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc/utils"
	"github.com/bitcoin-sv/go-broadcast-client/httpclient"
)

var ErrMissingTxID = errors.New("missing tx id")

func (a *ArcClient) QueryTransaction(ctx context.Context, txID string) (*broadcast.QueryTxResponse, error) {
	if a == nil {
		return nil, broadcast.ErrClientUndefined
	}

	result, err := queryTransaction(ctx, a, txID)
	if err != nil {
		return nil, arcutils.WithCause(errors.New("QueryTransaction: querying failed"), err)
	}

	err = validateQueryTxResponse(result)
	if err != nil {
		return nil, arcutils.WithCause(errors.New("QueryTransaction: validation of query tx response failed"), err)
	}

	a.Logger.Debug().Msgf("Got query tx response from miner: %s", result.Miner)
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

	if arc.headers != nil {
		pld.Headers = arc.headers
	}

	parseResponse := func(resp *http.Response) (*broadcast.QueryTxResponse, error) {
		return decodeQueryResponseBody(resp, arc)
	}

	return httpclient.RequestModel(
		ctx,
		arc.HTTPClient.DoRequest,
		pld,
		parseResponse,
		parseArcError,
	)
}

func decodeQueryResponseBody(resp *http.Response, arc *ArcClient) (*broadcast.QueryTxResponse, error) {
	base := broadcast.BaseTxResponse{}
	err := arcutils.DecodeResponseBody(resp.Body, &base)
	if err != nil {
		return nil, err
	}

	model := &broadcast.QueryTxResponse{
		BaseResponse: broadcast.BaseResponse{
			Miner: arc.apiURL,
		},
		BaseTxResponse: base,
	}

	return model, nil
}

func validateQueryTxResponse(model *broadcast.QueryTxResponse) error {
	if model.TxID == "" {
		return ErrMissingTxID
	}

	return nil
}
