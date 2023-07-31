package broadcast

import (
	"context"
	"fmt"
)

type BroadcastFactory func() Broadcaster

type Broadcaster interface {
	// BestQuoter
	// FastestQuoter
	// FeeQuoter
	// PolicyQuoter
	TransactionQuerier
	TransactionSubmitter
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

func (c *compositeBroadcaster) QueryTransaction(ctx context.Context, txID string) (*QueryTxResponse, error) {
	executionFuncs := make([]ExecutionFunc, len(c.broadcasters))
	for i, broadcaster := range c.broadcasters {
		executionFuncs[i] = func(ctx context.Context) (Result, error) {
			return broadcaster.QueryTransaction(ctx, txID)
		}
	}
	result, err := c.strategy.Execute(ctx, executionFuncs)
	if err != nil {
		return nil, err
	}

	// Convert result to QueryTxResponse
	queryTxResponse, ok := result.(*QueryTxResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return queryTxResponse, nil
}

func (c *compositeBroadcaster) SubmitTransaction(ctx context.Context, tx *Transaction) (*SubmitTxResponse, error) {
	executionFuncs := make([]ExecutionFunc, len(c.broadcasters))
	for i, broadcaster := range c.broadcasters {
		executionFuncs[i] = func(ctx context.Context) (Result, error) {
			return broadcaster.SubmitTransaction(ctx, tx)
		}
	}
	result, err := c.strategy.Execute(ctx, executionFuncs)
	if err != nil {
		return nil, err
	}

	// Convert result to QueryTxResponse
	submitTxResponse, ok := result.(*SubmitTxResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return submitTxResponse, nil
}
