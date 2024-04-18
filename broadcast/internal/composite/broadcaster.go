package composite

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

var DefaultStrategy = *OneByOne

type BroadcastFactory func() broadcast.Client

type compositeBroadcaster struct {
	broadcasters []broadcast.Client
	strategy     Strategy
}

func NewBroadcasterWithDefaultStrategy(factories ...BroadcastFactory) broadcast.Client {
	return NewBroadcaster(DefaultStrategy, factories...)
}

func NewBroadcaster(strategy Strategy, factories ...BroadcastFactory) broadcast.Client {
	var broadcasters []broadcast.Client
	for _, factory := range factories {
		broadcasters = append(broadcasters, factory())
	}
	return &compositeBroadcaster{
		broadcasters: broadcasters,
		strategy:     strategy,
	}
}

func (c *compositeBroadcaster) GetPolicyQuote(
	ctx context.Context,
) ([]*broadcast.PolicyQuoteResponse, *broadcast.SubmitFailure) {
	var policyQuotes []*broadcast.PolicyQuoteResponse

	for _, broadcaster := range c.broadcasters {
		singlePolicy, err := broadcaster.GetPolicyQuote(ctx)
		if err == nil {
			policyQuotes = append(policyQuotes, singlePolicy[0])
		}
	}

	if policyQuotes == nil {
		return nil, broadcast.Failure("GetPolicyQuote: ", broadcast.ErrNoMinerResponse)
	}

	return policyQuotes, nil
}

func (c *compositeBroadcaster) GetFeeQuote(ctx context.Context) ([]*broadcast.FeeQuote, *broadcast.SubmitFailure) {
	var feeQuotes []*broadcast.FeeQuote

	for _, broadcaster := range c.broadcasters {
		singleFeeQuote, err := broadcaster.GetFeeQuote(ctx)
		if err == nil {
			feeQuotes = append(feeQuotes, singleFeeQuote[0])
		}
	}

	if feeQuotes == nil {
		return nil, broadcast.Failure("GetFeeQuote: ", broadcast.ErrNoMinerResponse)
	}

	return feeQuotes, nil
}

func (c *compositeBroadcaster) QueryTransaction(
	ctx context.Context,
	txID string,
) (*broadcast.QueryTxResponse, *broadcast.SubmitFailure) {
	executionFuncs := make([]executionFunc, len(c.broadcasters))
	for i, broadcaster := range c.broadcasters {
		currentBroadcaster := broadcaster
		executionFuncs[i] = func(ctx context.Context) (Result, *broadcast.SubmitFailure) {
			return currentBroadcaster.QueryTransaction(ctx, txID)
		}
	}
	result, err := c.strategy.Execute(ctx, executionFuncs)
	if err != nil {
		return nil, err
	}

	// Convert result to QueryTxResponse
	queryTxResponse, ok := result.(*broadcast.QueryTxResponse)
	if !ok {
		return nil, broadcast.Failure(fmt.Sprintf("unexpected result type: %T", result), nil)
	}

	return queryTxResponse, nil
}

func (c *compositeBroadcaster) SubmitTransaction(
	ctx context.Context,
	tx *broadcast.Transaction,
	opts ...broadcast.TransactionOptFunc,
) (*broadcast.SubmitTxResponse, *broadcast.SubmitFailure) {
	executionFuncs := make([]executionFunc, len(c.broadcasters))
	for i, broadcaster := range c.broadcasters {
		currentBroadcaster := broadcaster
		executionFuncs[i] = func(ctx context.Context) (Result, *broadcast.SubmitFailure) {
			return currentBroadcaster.SubmitTransaction(ctx, tx)
		}
	}
	result, err := c.strategy.Execute(ctx, executionFuncs)
	if err != nil {
		return nil, err
	}

	// Convert result to SubmitTxResponse
	submitTxResponse, ok := result.(*broadcast.SubmitTxResponse)
	if !ok {
		return nil, broadcast.Failure(fmt.Sprintf("unexpected result type: %T", result), nil)
	}

	return submitTxResponse, nil
}

func (c *compositeBroadcaster) SubmitBatchTransactions(
	ctx context.Context,
	txs []*broadcast.Transaction,
	opts ...broadcast.TransactionOptFunc,
) (*broadcast.SubmitBatchTxResponse, *broadcast.SubmitFailure) {
	executionFuncs := make([]executionFunc, len(c.broadcasters))
	for i, broadcaster := range c.broadcasters {
		currentBroadcaster := broadcaster
		executionFuncs[i] = func(ctx context.Context) (Result, *broadcast.SubmitFailure) {
			return currentBroadcaster.SubmitBatchTransactions(ctx, txs)
		}
	}

	result, err := c.strategy.Execute(ctx, executionFuncs)
	if err != nil {
		return nil, err
	}

	// Convert result to []SubmitTxResponse
	submitTxResponse, ok := result.(*broadcast.SubmitBatchTxResponse)
	if !ok {
		return nil, broadcast.Failure(fmt.Sprintf("unexpected result type: %T", result), nil)
	}

	return submitTxResponse, nil
}
