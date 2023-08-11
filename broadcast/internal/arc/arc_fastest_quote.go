package arc

import (
	"context"
	"time"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

func (a *ArcClient) GetFastestQuote(
	ctx context.Context,
	timeout time.Duration,
) (*broadcast.FeeQuote, error) {
	if timeout.Seconds() == 0 {
		timeout = broadcast.DefaultFastestQuoteTimeout
	}
    
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	policyQuote, err := a.GetPolicyQuote(ctxWithTimeout)
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
