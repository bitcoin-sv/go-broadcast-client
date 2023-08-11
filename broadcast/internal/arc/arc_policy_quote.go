package arc

import (
	"context"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	arc_utils "github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc/utils"
	"github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/httpclient"
)

func (a *ArcClient) GetPolicyQuote(ctx context.Context) (*broadcast.PolicyQuoteResponse, error) {
	if a == nil {
		return nil, broadcast.ErrClientUndefined
	}

	model, err := getPolicyQuote(ctx, a)
	if err != nil {
		return nil, err
	}

	model.Miner = a.apiURL

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

	resp, err := arc.HTTPClient.DoRequest(ctx, pld)
	if err != nil {
		return nil, arc_utils.HandleHttpError(err)
	}

	model := &broadcast.PolicyQuoteResponse{}
	err = arc_utils.DecodeResponseBody(resp.Body, model)
	if err != nil {
		return nil, err
	}
	return model, nil
}
