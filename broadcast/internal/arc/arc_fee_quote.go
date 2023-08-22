package arc

import (
	"context"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

func (a *ArcClient) GetFeeQuote(ctx context.Context) ([]*broadcast.FeeQuote, error) {
	policyQuotes, err := a.GetPolicyQuote(ctx)
	if err != nil {
		return nil, err
	}

	feeQuote := &broadcast.FeeQuote{
		BaseResponse: broadcast.BaseResponse{Miner: a.apiURL},
		MiningFee:    policyQuotes[0].Policy.MiningFee,
		Timestamp:    policyQuotes[0].Timestamp,
	}

	feeQuotes := []*broadcast.FeeQuote{feeQuote}

	return feeQuotes, nil
}
