package arc

import (
	"context"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

func (a *ArcClient) GetFeeQuote(ctx context.Context) ([]*broadcast.FeeQuote, *broadcast.FailureResponse) {
	if a == nil {
		return nil, broadcast.Failure("GetFeeQuote:", broadcast.ErrClientUndefined)
	}

	policyQuotes, err := a.GetPolicyQuote(ctx)
	if err != nil {
		return nil, broadcast.Failure("GetFeeQuote: request failed", err)
	}

	feeQuote := &broadcast.FeeQuote{
		BaseResponse: broadcast.BaseResponse{Miner: a.apiURL},
		MiningFee:    policyQuotes[0].Policy.MiningFee,
		Timestamp:    policyQuotes[0].Timestamp,
	}

	feeQuotes := []*broadcast.FeeQuote{feeQuote}

	a.Logger.Debug().Msgf("Got fee quote from miner: %s", feeQuote.Miner)
	return feeQuotes, nil
}
