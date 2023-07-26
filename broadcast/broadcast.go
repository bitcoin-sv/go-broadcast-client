package broadcast

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/go-broadcast-client/models"
)

type BroadcastFactory func() Broadcaster

type Broadcaster interface {
	// BestQuoter
	// FastestQuoter
	// FeeQuoter
	// PolicyQuoter
	TransactionQuerier
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

func (c *compositeBroadcaster) QueryTransaction(ctx context.Context, txID string) (*models.QueryTxResponse, error) {
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
	queryTxResponse, ok := result.(*models.QueryTxResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return queryTxResponse, nil
}
