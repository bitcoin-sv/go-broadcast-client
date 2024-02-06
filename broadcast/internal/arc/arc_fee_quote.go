package arc

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	arc_utils "github.com/bitcoin-sv/go-broadcast-client/broadcast/internal/arc/utils"
)

func (a *ArcClient) GetFeeQuote(ctx context.Context) ([]*broadcast.FeeQuote, error) {
	policyQuotes, err := a.GetPolicyQuote(ctx)
	if err != nil {
		return nil, arc_utils.WithCause(errors.New("GetFeeQuote: request failed"), err)
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
