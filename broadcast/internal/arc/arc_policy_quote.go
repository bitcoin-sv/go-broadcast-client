package arc

import (
	"context"
	"net/http"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	arc_utils "github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc/utils"
	"github.com/bitcoin-sv/go-broadcast-client/httpclient"
)

func (a *ArcClient) GetPolicyQuote(ctx context.Context) ([]*broadcast.PolicyQuoteResponse, error) {
	if a == nil {
		return nil, broadcast.ErrClientUndefined
	}

	model, err := getPolicyQuote(ctx, a)
	if err != nil {
		a.Logger.Error().Msgf("Failed to get policy quote: %s", err.Error())
		return nil, err
	}

	model.Miner = a.apiURL
	models := []*broadcast.PolicyQuoteResponse{model}

	a.Logger.Debug().Msgf("Got policy quote from miner: %s", model.Miner)
	return models, nil
}

func decodePolicyQuoteResponseBody(resp *http.Response) (*broadcast.PolicyQuoteResponse, error) {
	model := &broadcast.PolicyQuoteResponse{}
	err := arc_utils.DecodeResponseBody(resp.Body, model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func getPolicyQuote(ctx context.Context, arc *ArcClient) (*broadcast.PolicyQuoteResponse, error) {
	url := arc.apiURL + arcPolicyQuoteRoute
	pld := httpclient.NewPayload(
		httpclient.GET,
		url,
		arc.token,
		nil,
	)

	if arc.deploymentID != "" {
		pld.AddHeader(XDeploymentIDHeader, arc.deploymentID)
	}

	return httpclient.RequestModel(
		ctx,
		arc.HTTPClient.DoRequest,
		pld,
		decodePolicyQuoteResponseBody,
		parseArcError,
	)
}
