package composite

import (
	"context"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

// StrategyName is the name of the strategy.
type StrategyName string

const (
	// OneByOneStrategy is the name of the one by one strategy.
	OneByOneStrategy StrategyName = "OneByOneStrategy"
)

// Result is the interface for the result of an execution func.
type Result interface{}

type executionFunc func(context.Context) (Result, error)

// StrategyExecutionFunc is the function that executes the strategy.
type StrategyExecutionFunc func(context.Context, []executionFunc) (Result, error)

// Strategy is the composite strategy.
type Strategy struct {
	name          StrategyName
	executionFunc StrategyExecutionFunc
}

// New creates a new composite strategy.
func New(name StrategyName) (*Strategy, error) {
	switch name {
	case OneByOneStrategy:
		return &Strategy{name: name, executionFunc: OneByOne.executionFunc}, nil
	default:
		return nil, broadcast.ErrStrategyUnkown
	}
}

// Execute executes the strategy.
func (s *Strategy) Execute(ctx context.Context, executionFuncs []executionFunc) (Result, error) {
	return s.executionFunc(ctx, executionFuncs)
}

var (
	OneByOne = &Strategy{name: OneByOneStrategy, executionFunc: func(ctx context.Context, executionFuncs []executionFunc) (Result, error) {
		for _, executionFunc := range executionFuncs {
			result, err := executionFunc(ctx)
			if err != nil {
				continue
			}
			return result, nil
		}
		return nil, broadcast.ErrAllBroadcastersFailed
	}}
)
