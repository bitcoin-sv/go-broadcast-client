package arc

import (
	"context"
	"errors"
	"net/http"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	arc_utils "github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc/utils"

	"github.com/bitcoin-sv/go-broadcast-client/httpclient"
)

var ErrMissingTxID = errors.New("missing tx id")

func (a *ArcClient) QueryTransaction(ctx context.Context, txID string) (*broadcast.QueryTxResponse, broadcast.ArcFailure) {
	if a == nil {
		return nil, broadcast.Failure("QueryTransaction:", broadcast.ErrClientUndefined)
	}

	result, err := queryTransaction(ctx, a, txID)
	if err != nil {
		return nil, broadcast.Failure("QueryTransaction: querying failed", err)
	}

	err = validateQueryTxResponse(result)
	if err != nil {
		return nil, broadcast.Failure("QueryTransaction: validation of query tx response failed", err)
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

	if arc.deploymentID != "" {
		pld.AddHeader(XDeploymentIDHeader, arc.deploymentID)
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
	err := arc_utils.DecodeResponseBody(resp.Body, &base)
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
