package arc

import (
	"context"
	"encoding/json"
	"io"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
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
		return nil, err
	}

	model, err := decodeQueryPolicyBody(resp.Body)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func decodeQueryPolicyBody(body io.ReadCloser) (*broadcast.PolicyQuoteResponse, error) {
	model := broadcast.PolicyQuoteResponse{}
	err := json.NewDecoder(body).Decode(&model)
	if err != nil {
		return nil, broadcast.ErrUnableToDecodeResponse
	}

	return &model, nil
}
