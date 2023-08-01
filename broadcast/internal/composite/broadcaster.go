package composite

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-api"
)

var DefaultStrategy = *OneByOne

type BroadcastFactory func() broadcast_api.Client

type compositeBroadcaster struct {
	broadcasters []broadcast_api.Client
	strategy     Strategy
}

func NewBroadcasterWithDefaultStrategy(factories ...BroadcastFactory) broadcast_api.Client {
	return NewBroadcaster(DefaultStrategy, factories...)
}

func NewBroadcaster(strategy Strategy, factories ...BroadcastFactory) broadcast_api.Client {
	var broadcasters []broadcast_api.Client
	for _, factory := range factories {
		broadcasters = append(broadcasters, factory())
	}
	return &compositeBroadcaster{
		broadcasters: broadcasters,
		strategy:     strategy,
	}
}

func (c *compositeBroadcaster) QueryTransaction(ctx context.Context, txID string) (*broadcast_api.QueryTxResponse, error) {
	executionFuncs := make([]executionFunc, len(c.broadcasters))
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
	queryTxResponse, ok := result.(*broadcast_api.QueryTxResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return queryTxResponse, nil
}

func (c *compositeBroadcaster) SubmitTransaction(ctx context.Context, tx *broadcast_api.Transaction) (*broadcast_api.SubmitTxResponse, error) {
	executionFuncs := make([]executionFunc, len(c.broadcasters))
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
	submitTxResponse, ok := result.(*broadcast_api.SubmitTxResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return submitTxResponse, nil
}
