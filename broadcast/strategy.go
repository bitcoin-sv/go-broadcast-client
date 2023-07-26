package broadcast

import (
	"context"

	errors "github.com/bitcoin-sv/go-broadcast-client"
)

type StrategyName string

const (
	OneByOneStrategy StrategyName = "OneByOneStrategy"
)

type Result interface{}

type ExecutionFunc func(context.Context) (Result, error)
type StrategyExecutionFunc func(context.Context, []ExecutionFunc) (Result, error)

type Strategy struct {
	name          StrategyName
	executionFunc StrategyExecutionFunc
}

func New(name StrategyName, executionFunc StrategyExecutionFunc) *Strategy {
	return &Strategy{name: name, executionFunc: executionFunc}
}

func (s *Strategy) Execute(ctx context.Context, executionFuncs []ExecutionFunc) (Result, error) {
	return s.executionFunc(ctx, executionFuncs)
}

var (
	OneByOne = New(OneByOneStrategy, func(ctx context.Context, executionFuncs []ExecutionFunc) (Result, error) {
		for _, executionFunc := range executionFuncs {
			result, err := executionFunc(ctx)
			if err != nil {
				continue
			}
			return result, nil
		}
		return nil, errors.ErrAllBroadcastersFailed
	})
)
