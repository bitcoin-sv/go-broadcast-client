package composite

import (
	"context"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
)

type StrategyName string

const (
	// OneByOneStrategy is a strategy that executes the execution funcs one by one until one of them succeeds.
	// If all execution funcs fail, then the strategy returns an error.
	// The error is: ErrAllBroadcastersFailed.
	OneByOneStrategy StrategyName = "OneByOneStrategy"
)

type Result interface{}

type executionFunc func(context.Context) (Result, broadcast.ArcFailure)

type StrategyExecutionFunc func(context.Context, []executionFunc) (Result, broadcast.ArcFailure)

// Strategy is a component designed to offer flexibility in selecting a communication approach
// for interacting with multiple broadcasting services, such as multiple Arc services.
type Strategy struct {
	name          StrategyName
	executionFunc StrategyExecutionFunc
}

func New(name StrategyName) (*Strategy, error) {
	switch name {
	case OneByOneStrategy:
		return &Strategy{name: name, executionFunc: OneByOne.executionFunc}, nil
	default:
		return nil, broadcast.ErrStrategyUnkown
	}
}

func (s *Strategy) Execute(ctx context.Context, executionFuncs []executionFunc) (Result, broadcast.ArcFailure) {
	return s.executionFunc(ctx, executionFuncs)
}

var (
	OneByOne = &Strategy{name: OneByOneStrategy, executionFunc: func(ctx context.Context, executionFuncs []executionFunc) (Result, broadcast.ArcFailure) {
		for _, executionFunc := range executionFuncs {
			result, err := executionFunc(ctx)
			if err != nil {
				continue
			}
			return result, nil
		}
		return nil, broadcast.Failure("", broadcast.ErrAllBroadcastersFailed)
	}}
)
