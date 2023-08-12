package composite

import (
	"context"
	"fmt"
	"sync"
	"time"

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

func (c *compositeBroadcaster) GetPolicyQuote(ctx context.Context) (*broadcast.PolicyQuoteResponse, error) {
	executionFuncs := make([]executionFunc, len(c.broadcasters))
	for i, broadcaster := range c.broadcasters {
		executionFuncs[i] = func(ctx context.Context) (Result, error) {
			return broadcaster.GetPolicyQuote(ctx)
		}
	}

	result, err := c.strategy.Execute(ctx, executionFuncs)
	if err != nil {
		return nil, err
	}

	policyQuoteResponse, ok := result.(*broadcast.PolicyQuoteResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return policyQuoteResponse, nil
}

func (c *compositeBroadcaster) GetFeeQuote(ctx context.Context) (*broadcast.FeeQuote, error) {
	executionFuncs := make([]executionFunc, len(c.broadcasters))
	for i, broadcaster := range c.broadcasters {
		executionFuncs[i] = func(ctx context.Context) (Result, error) {
			return broadcaster.GetFeeQuote(ctx)
		}
	}

	result, err := c.strategy.Execute(ctx, executionFuncs)
	if err != nil {
		return nil, err
	}

	feeQuote, ok := result.(*broadcast.FeeQuote)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return feeQuote, nil
}

func (c *compositeBroadcaster) GetBestQuote(ctx context.Context) (*broadcast.FeeQuote, error) {
	fees := make(chan *broadcast.FeeQuote, len(c.broadcasters))
	var wg sync.WaitGroup

	for _, broadcaster := range c.broadcasters {
		wg.Add(1)
		go func(ctx context.Context, b broadcast.Client) {
			defer wg.Done()
			feeQuote, err := b.GetFeeQuote(ctx)
			if err == nil {
				fees <- feeQuote
			}
		}(ctx, broadcaster)
	}

	wg.Wait()
	close(fees)

	var bestQuote *broadcast.FeeQuote = nil
	for fee := range fees {
		if bestQuote == nil {
			bestQuote = fee
		} else {
			feePer1000Bytes := (1000 / fee.MiningFee.Bytes) * fee.MiningFee.Satoshis
			bestFeePer1000Bytes := (1000 / bestQuote.MiningFee.Bytes) * bestQuote.MiningFee.Satoshis

			if feePer1000Bytes < bestFeePer1000Bytes {
				bestQuote = fee
			}
		}
	}

	if bestQuote == nil {
		return nil, broadcast.ErrNoMinerResponse
	}

	return bestQuote, nil
}

func (c *compositeBroadcaster) GetFastestQuote(ctx context.Context, timeout time.Duration) (*broadcast.FeeQuote, error) {
	if timeout.Seconds() == 0 {
		timeout = broadcast.DefaultFastestQuoteTimeout
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	feeChannel := make(chan *broadcast.FeeQuote, len(c.broadcasters))

	var wg sync.WaitGroup
	for _, broadcaster := range c.broadcasters {
		wg.Add(1)
		go func(ctxWithTimeout context.Context, broadcaster broadcast.Client) {
			defer wg.Done()
			feeQuote, err := broadcaster.GetFeeQuote(ctxWithTimeout)
			if err == nil {
				feeChannel <- feeQuote
			}
		}(ctxWithTimeout, broadcaster)
	}

	go func() {
		wg.Wait()
		close(feeChannel)
	}()

	fastQuote := <-feeChannel

	if fastQuote == nil {
		return nil, broadcast.ErrNoMinerResponse
	}

	return fastQuote, nil
}

func (c *compositeBroadcaster) QueryTransaction(ctx context.Context, txID string) (*broadcast.QueryTxResponse, error) {
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
	queryTxResponse, ok := result.(*broadcast.QueryTxResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return queryTxResponse, nil
}

func (c *compositeBroadcaster) SubmitTransaction(ctx context.Context, tx *broadcast.Transaction) (*broadcast.SubmitTxResponse, error) {
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

	// Convert result to SubmitTxResponse
	submitTxResponse, ok := result.(*broadcast.SubmitTxResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return submitTxResponse, nil
}

func (c *compositeBroadcaster) SubmitBatchTransactions(ctx context.Context, txs []*broadcast.Transaction) ([]*broadcast.SubmitTxResponse, error) {
	executionFuncs := make([]executionFunc, len(c.broadcasters))
	for i, broadcaster := range c.broadcasters {
		executionFuncs[i] = func(ctx context.Context) (Result, error) {
			return broadcaster.SubmitBatchTransactions(ctx, txs)
		}
	}

	result, err := c.strategy.Execute(ctx, executionFuncs)
	if err != nil {
		return nil, err
	}

	// Convert result to []SubmitTxResponse
	submitTxResponse, ok := result.([]*broadcast.SubmitTxResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	return submitTxResponse, nil
}
