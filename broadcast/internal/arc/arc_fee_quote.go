package arc

import (
	"context"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

func (a *ArcClient) GetFeeQuote(ctx context.Context) (*broadcast.FeeQuote, error) {
	policyQuote, err := a.GetPolicyQuote(ctx)
	if err != nil {
		return nil, err
	}

	feeQuote := &broadcast.FeeQuote{
		Miner:     a.apiURL,
		MiningFee: policyQuote.Policy.MiningFee,
		Timestamp: policyQuote.Timestamp,
	}
	return feeQuote, nil
}
