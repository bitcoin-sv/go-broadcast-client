package broadcast

import (
	"context"
)

type BroadcastFactory func() Broadcaster

type Broadcaster interface {
	BestQuoter
	// FastestQuoter
	// FeeQuoter
	// PolicyQuoter
	// TransactionQuerier
	// TransactionSubmitter
	// TransactionsSubmitter
}

type compositeBroadcaster struct {
	broadcasters []Broadcaster
	strategy     Strategy
}

func NewCompositeBroadcaster(strategy Strategy, factories ...BroadcastFactory) Broadcaster {
	var broadcasters []Broadcaster
	for _, factory := range factories {
		broadcasters = append(broadcasters, factory())
	}
	return &compositeBroadcaster{
		broadcasters: broadcasters,
		strategy:     strategy,
	}
}

func (c *compositeBroadcaster) BestQuote(ctx context.Context, feeCategory, feeType string) error {
	executionFuncs := make([]func(context.Context) error, len(c.broadcasters))
	for i, broadcaster := range c.broadcasters {
		executionFuncs[i] = func(ctx context.Context) error {
			return broadcaster.BestQuote(ctx, feeCategory, feeType)
		}
	}
	return c.strategy.Execute(ctx, executionFuncs)
}
